package user

import (
	"context"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/validator"
)

type UserService interface {
	Register(ctx context.Context, user models.RegisterRequest) (string, error)
	Login(ctx context.Context, user models.LoginRequest) (string, error)
}

type UserHandler struct {
	service   UserService
	validator *validator.Validator
	config    config.Config
}

const logContext = "UserHandler"

func New(cfg config.Config, us UserService, validator *validator.Validator) *UserHandler {
	return &UserHandler{
		service:   us,
		validator: validator,
		config:    cfg,
	}
}
