package api

import (
	"blackwallgroup/queue"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Users(c *gin.Context)
	Withdraw(c *gin.Context)
	Deposit(c *gin.Context)
}

type Server struct {
	engine           *gin.Engine
	postgresPool     *pgxpool.Pool
	transactionQueue *queue.Queue
	userService      UserService
}

func NewServer(postgresPool *pgxpool.Pool, transactionQueue *queue.Queue, userService UserService) *Server {
	engine := gin.Default()
	engine.Use(gin.Recovery())

	return &Server{
		engine:           engine,
		postgresPool:     postgresPool,
		userService:      userService,
		transactionQueue: transactionQueue,
	}
}

func (s *Server) Start(addr string) error {
	api := s.engine.Group("/api")
	api.POST("/users", s.userService.Users)
	api.POST("/withdraw", s.userService.Withdraw)
	api.POST("/deposit", s.userService.Deposit)

	return s.engine.Run(addr)
}
