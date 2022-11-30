package repository

import (
	"blackwallgroup/internal/model"
	"blackwallgroup/queue"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type Storage struct {
	pool             *pgxpool.Pool
	transactionQueue *queue.Queue
}

func NewStorage(postgresPool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: postgresPool,
	}
}

func (s *Storage) Users(c *gin.Context) {
	sql, _, _ := squirrel.Select("*").From("clients.clients").ToSql()

	rows, err := s.pool.Query(c.Request.Context(), sql)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "database query error"})
		return
	}

	var users []model.User
	err = pgxscan.NewRowScanner(rows).Scan(&users)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "scan rows error"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s *Storage) Withdraw(c *gin.Context) {
	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	s.transactionQueue.Put(req)
	s.transactionQueue.Wait(req)

	tx, err := s.pool.BeginTx(c.Request.Context(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error begin transaction"})
		return
	}
	defer tx.Rollback(c.Request.Context())

	sql, _, _ := squirrel.Select("balance").From("clients.clients").Where("id = ?", req.ID).ToSql()

	var user model.User
	err = tx.QueryRow(c.Request.Context(), sql).Scan(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	newBalance := user.Balance - req.Amount
	if newBalance < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "small balance error"})
		return
	}

	user.Balance = newBalance

	sql, _, _ = squirrel.Update("clients").Set("balance = ?", newBalance).ToSql()
	_, err = tx.Exec(c.Request.Context(), sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	err = tx.Commit(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error commit transaction"})
		return
	}

	s.transactionQueue.Release(req)

	c.JSON(http.StatusOK, user)
}

func (s *Storage) Deposit(c *gin.Context) {
	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	sql, _, _ := squirrel.Select("balance").From("clients.clients").Where("id = ?", req.ID).ToSql()

	tx, err := s.pool.BeginTx(c.Request.Context(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error begin transaction"})
		return
	}
	defer tx.Rollback(c.Request.Context())

	var user model.User
	err = tx.QueryRow(c.Request.Context(), sql).Scan(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	user.Balance = user.Balance + req.Amount

	sql, _, _ = squirrel.Update("clients").Set("balance = ?", user.Balance).ToSql()

	_, err = tx.Exec(c.Request.Context(), sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	err = tx.Commit(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error commit transaction"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}
