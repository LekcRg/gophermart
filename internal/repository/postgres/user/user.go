package user

import (
	"context"
	"fmt"

	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, db *pgxpool.Pool) *UserPostgres {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(20) UNIQUE NOT NULL,
		passhash varchar(72) NOT NULL
	)`
	_, err := db.Exec(ctx, query)
	if err != nil {
		logger.Log.Error("Create user table error",
			zap.Error(err))
	}
	return &UserPostgres{
		db: db,
	}
}

func (up *UserPostgres) Create(ctx context.Context, user models.User) error {
	fmt.Printf("%+v\n", user)
	query := `INSERT INTO users (login, passhash)
	VALUES ($1, $2)`
	_, err := up.db.Exec(ctx, query, user.Login, user.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}
