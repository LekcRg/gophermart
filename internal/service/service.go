package service

import (
	"context"

	"github.com/LekcRg/gophermart/internal/accrual"
	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/service/orders"
	"github.com/LekcRg/gophermart/internal/service/user"
	"github.com/LekcRg/gophermart/internal/service/withdraw"
	"github.com/LekcRg/gophermart/internal/validator"
)

type Service struct {
	User     user.UserService
	Orders   orders.OrdersService
	Withdraw withdraw.WithdrawService
}

func New(
	ctx context.Context,
	db *repository.Repository, validator *validator.Validator,
	cfg config.Config, req *accrual.Accrual,
) *Service {
	return &Service{
		User:     *user.New(db.User, validator, cfg),
		Orders:   *orders.New(ctx, db.Orders, validator, cfg, req, db.User),
		Withdraw: *withdraw.New(db.Withdraw, db.User, validator),
	}
}
