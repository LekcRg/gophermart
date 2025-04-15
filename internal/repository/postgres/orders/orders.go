package orders

import (
	"context"

	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type OrdersPostgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, db *pgxpool.Pool) *OrdersPostgres {
	// create orders table
	query := `DO $$
	BEGIN
		CREATE TYPE order_status AS ENUM ( 'NEW',
			'PROCESSING',
			'INVALID',
			'PROCESSED'
	);
	EXCEPTION
	WHEN duplicate_object THEN
		NULL;
	END
	$$;

	CREATE TABLE IF NOT EXISTS orders (
		id serial NOT NULL PRIMARY KEY,
		number varchar(50) UNIQUE,
		status order_status NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP NOT NULL DEFAULT now(),
		user_login varchar(30) NOT NULL REFERENCES users (LOGIN)
	);`

	_, err := db.Exec(ctx, query)
	if err != nil {
		logger.Log.Error("Error while create orders table",
			zap.Error(err))
	}

	return &OrdersPostgres{
		db: db,
	}
}

func (op OrdersPostgres) Create(
	ctx context.Context, num string, status string, user models.DBUser,
) error {
	query := `INSERT INTO orders (number, status, user_login, updated_at)
		values($1, $2, $3, now())`

	_, err := op.db.Exec(ctx, query, num, status, user.Login)
	if err != nil {
		return err
	}
	return nil
}
