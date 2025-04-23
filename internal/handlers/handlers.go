package handlers

import (
	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/handlers/orders"
	"github.com/LekcRg/gophermart/internal/handlers/user"
	"github.com/LekcRg/gophermart/internal/handlers/withdraw"
	"github.com/LekcRg/gophermart/internal/service"
	"github.com/LekcRg/gophermart/internal/validator"
)

type Handlers struct {
	User     *user.UserHandler
	Orders   *orders.OrdersHandler
	Withdraw *withdraw.WithdrawHandler
}

func New(cfg config.Config, s *service.Service, validator *validator.Validator) *Handlers {
	return &Handlers{
		User:     user.New(cfg, &s.User, validator),
		Orders:   orders.New(cfg, &s.Orders, validator),
		Withdraw: withdraw.New(&s.Withdraw, validator),
	}
}
