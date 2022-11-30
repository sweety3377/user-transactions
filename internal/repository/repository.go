package repository

import (
	"blackwallgroup/internal/model"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(postgresPool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: postgresPool,
	}
}

func (s *Storage) Users(c *gin.Context) {
	sql, _, _ := squirrel.Select("*").From("client").ToSql()

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

	sql, _, _ := squirrel.Select("balance").Where("id = ?", req.ID).ToSql()

	var user model.User
	err := s.pool.QueryRow(c.Request.Context(), sql).Scan(&user)
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
	_, err = s.pool.Exec(c.Request.Context(), sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Storage) Deposit(c *gin.Context) {
	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	sql, _, _ := squirrel.Select("balance").Where("id = ?", req.ID).ToSql()

	var user model.User
	err := s.pool.QueryRow(c.Request.Context(), sql).Scan(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	user.Balance = user.Balance + req.Amount

	sql, _, _ = squirrel.Update("clients").Set("balance = ?", user.Balance).ToSql()

	_, err = s.pool.Exec(c.Request.Context(), sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind user model error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
