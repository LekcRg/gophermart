package withdraw

import (
	"context"

	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type WithdrawPostgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, db *pgxpool.Pool) *WithdrawPostgres {
	query := `CREATE TABLE IF NOT EXISTS withdraw (
		id SERIAL PRIMARY KEY,
		order_id varchar(50) NOT NULL,
		sum DOUBLE PRECISION DEFAULT 0,
		user_login varchar(30) NOT NULL REFERENCES users (login),
		processed_at TIMESTAMP NOT NULL DEFAULT now()
	)`
	_, err := db.Exec(ctx, query)
	if err != nil {
		logger.Log.Error("create withdraw table error",
			zap.Error(err))
	}

	return &WithdrawPostgres{
		db: db,
	}
}

func (wp *WithdrawPostgres) CreateWithdraw(ctx context.Context, userLogin string, withdraw models.WithdrawRequest) error {
	query := `INSERT INTO withdraw (sum, order_id, user_login)
		values($1, $2, $3)`
	_, err := wp.db.Exec(ctx, query, withdraw.Sum, withdraw.Order, userLogin)
	if err != nil {
		return err
	}

	// TODO: ошибки заказа, что его нету и т.п.

	return nil
}

func (wp *WithdrawPostgres) GetByUserLogin(
	ctx context.Context, userLogin string,
) (models.WithdrawList, error) {
	list := models.WithdrawList{}

	query := `SELECT SUM, order_id, processed_at FROM withdraw 
		WHERE user_login = $1`
	rows, err := wp.db.Query(ctx, query, userLogin)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		withdraw := models.Withdraw{}
		err := rows.Scan(&withdraw.Sum, &withdraw.Order, &withdraw.ProcessedAt)
		if err != nil {
			return list, err
		}

		list = append(list, withdraw)
	}

	return list, nil
}
