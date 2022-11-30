package service

import (
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	Users(c *gin.Context)
	Withdraw(c *gin.Context)
	Deposit(c *gin.Context)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (u *UserService) Users(c *gin.Context) {
	u.repository.Users(c)
}

func (u *UserService) Withdraw(c *gin.Context) {
	u.repository.Withdraw(c)
}

func (u *UserService) Deposit(c *gin.Context) {
	u.repository.Deposit(c)
}
