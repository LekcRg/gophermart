package service

import (
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/service/user"
)

type Service struct {
	User user.UserService
}

func New(db *repository.Repository) *Service {
	return &Service{
		User: *user.New(db.User),
	}
}
