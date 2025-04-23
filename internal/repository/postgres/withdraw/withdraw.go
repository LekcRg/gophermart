package withdraw

import (
	"context"

	"github.com/LekcRg/gophermart/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WithdrawPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *WithdrawPostgres {
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
