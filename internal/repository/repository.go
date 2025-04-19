package repository

import (
	"context"
	"fmt"

	"github.com/LekcRg/gophermart/internal/models"
)

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

var (
	ErrOrdersRegisteredOtherUser = fmt.Errorf(
		"order has already been uploaded by another user")
	ErrOrdersRegisteredThisUser = fmt.Errorf(
		"order has already been uploaded by this user")
)

type Repository struct {
	User   UserRepository
	Orders OrdersRepository
}

type RepositoryProvider interface {
	GetRepositories() *Repository
	Close()
}

type UserRepository interface {
	Create(context.Context, models.DBUser) (*models.DBUser, error)
	Login(context.Context, models.LoginRequest) (*models.DBUser, error)
	UpdateBalance(ctx context.Context, userLogin string, balance float64) error
	GetBalance(ctx context.Context, userLogin string) (models.UserBalance, error)
}

type OrdersRepository interface {
	GetOrdersByUserLogin(ctx context.Context, userLogin string) ([]models.OrderDB, error)
	Create(ctx context.Context, order models.OrderCreateDB, user models.DBUser) error
	UpdateOrder(ctx context.Context, orderID, status string, accrual float64) error
}
