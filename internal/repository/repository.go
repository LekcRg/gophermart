package repository

import (
	"context"

	"github.com/LekcRg/gophermart/internal/models"
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
}

type OrdersRepository interface {
	// GetOne() error
	Create(ctx context.Context, order string, status string, user models.DBUser) error
}

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)
