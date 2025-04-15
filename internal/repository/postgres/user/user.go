package user

import (
	"context"

	"github.com/LekcRg/gophermart/internal/errs"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, db *pgxpool.Pool) *UserPostgres {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(30) UNIQUE NOT NULL,
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

func (up *UserPostgres) Create(
	ctx context.Context, user models.DBUser,
) (*models.DBUser, error) {
	query := `INSERT INTO users (login, passhash) VALUES ($1, $2)
	RETURNING id, login`
	row := up.db.QueryRow(ctx, query, user.Login, user.PasswordHash)
	var userDB models.DBUser
	err := row.Scan(&userDB.ID, &userDB.Login)
	if err != nil {
		return nil, err
	}

	return &userDB, nil
}

func (up *UserPostgres) Login(
	ctx context.Context, user models.LoginRequest,
) (*models.DBUser, error) {
	query := `SELECT id, login, passhash FROM users WHERE login=$1`
	row := up.db.QueryRow(ctx, query, user.Login)
	var userDB models.DBUser
	err := row.Scan(&userDB.ID, &userDB.Login, &userDB.PasswordHash)
	if err != nil {
		return nil, err
	}
	if userDB.PasswordHash == "" {
		return nil, errs.ErrNotFoundUser
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash), []byte(user.Password))
	if err != nil {
		// maybe ErrNotFoundUser error
		return nil, errs.ErrIncorrectPassword
	}
	userDB.PasswordHash = ""

	return &userDB, nil
}
