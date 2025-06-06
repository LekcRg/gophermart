package user

import (
	"context"
	"database/sql"

	"github.com/LekcRg/gophermart/internal/errs"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *UserPostgres {
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

func (up *UserPostgres) UpdateBalance(
	ctx context.Context, userLogin string, balance float64,
) error {
	query := `UPDATE users SET balance = balance + $1 WHERE login = $2;`
	_, err := up.db.Exec(ctx, query, balance, userLogin)
	if err != nil {
		return err
	}

	return nil
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

func (up *UserPostgres) GetBalance(
	ctx context.Context, userLogin string,
) (models.UserBalance, error) {
	query := `SELECT balance, withdrawn FROM users WHERE login = $1`
	user := models.UserBalance{}
	row := up.db.QueryRow(ctx, query, userLogin)

	err := row.Scan(&user.Balance, &user.Withdrawn)
	if err != nil {
		return user, nil
	}

	return user, nil
}

func (up *UserPostgres) WithdrawBalance(
	ctx context.Context, UserLogin string, withdraw models.WithdrawRequest,
) error {
	query := `UPDATE users SET balance = balance - $1, withdrawn = withdrawn + $1
	WHERE login = $2 AND balance >= $1
	RETURNING users.balance`
	row := up.db.QueryRow(ctx, query, withdraw.Sum, UserLogin)

	var newBalance sql.NullFloat64
	err := row.Scan(&newBalance)
	if err != nil && err == pgx.ErrNoRows {
		return errs.ErrUserSmallBalance
	} else if err != nil {
		return err
	}

	return nil
}
