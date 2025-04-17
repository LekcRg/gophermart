package service

import (
	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/request"
	"github.com/LekcRg/gophermart/internal/service/orders"
	"github.com/LekcRg/gophermart/internal/service/user"
	"github.com/LekcRg/gophermart/internal/validator"
)

type Service struct {
	User   user.UserService
	Orders orders.OrdersService
}

func New(
	db *repository.Repository, validator *validator.Validator,
	cfg config.Config, req *request.Request,
) *Service {
	return &Service{
		User:   *user.New(db.User, validator, cfg),
		Orders: *orders.New(db.Orders, validator, cfg, req),
	}
}
