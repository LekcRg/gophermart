package service

import (
	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/service/user"
	"github.com/LekcRg/gophermart/internal/validator"
)

type Service struct {
	User user.UserService
}

func New(
	db *repository.Repository, validator *validator.Validator,
	cfg config.Config,
) *Service {
	return &Service{
		User: *user.New(db.User, validator, cfg),
	}
}
