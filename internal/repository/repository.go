package repository

import (
	"context"

	"github.com/LekcRg/gophermart/internal/models"
)

type Repository struct {
	User UserRepository
}

type RepositoryProvider interface {
	GetRepositories() *Repository
	Close()
}

type UserRepository interface {
	Create(context.Context, models.User) error
}
