package orders

import (
	"context"

	"github.com/LekcRg/gophermart/internal/errs"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/pgutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type OrdersPostgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, db *pgxpool.Pool) *OrdersPostgres {
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
		order_id varchar(50) NOT NULL UNIQUE,
		status order_status NOT NULL,
		accrual double precision,
		uploaded_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP NOT NULL DEFAULT now(),
		user_login varchar(30) NOT NULL REFERENCES users (LOGIN)
	);
`

	_, err := db.Exec(ctx, query)
	if err != nil {
		logger.Log.Error("Error while create orders table",
			zap.Error(err))
	}

	return &OrdersPostgres{
		db: db,
	}
}

func (op OrdersPostgres) GetOrdersByUserLogin(
	ctx context.Context, userLogin string,
) ([]models.OrderDB, error) {
	query := `select order_id, status, accrual, user_login, uploaded_at from orders
	where user_login = $1`

	orders := []models.OrderDB{}
	rows, err := op.db.Query(ctx, query, userLogin)
	if err != nil {
		return orders, err
	}

	for rows.Next() {
		var order models.OrderDB

		err := rows.Scan(
			&order.OrderID, &order.Status, &order.Accrual,
			&order.UserLogin, &order.UploadedAt)
		if err != nil {
			return orders, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (op OrdersPostgres) GetOrderByOrderID(
	ctx context.Context, orderID string,
) (models.OrderDB, error) {
	query := `select id, order_id, status, user_login, uploaded_at, updated_at 
	from orders where order_id = $1`

	order := models.OrderDB{}
	row := op.db.QueryRow(ctx, query, orderID)
	err := row.Scan(
		&order.ID, &order.OrderID, &order.Status,
		&order.UserLogin, &order.UploadedAt, &order.UpdatedAt)
	if err != nil {
		return models.OrderDB{}, err
	}

	return order, nil
}

func (op OrdersPostgres) UpdateOrder(
	ctx context.Context, orderID, status string, accrual float64,
) error {
	query := `UPDATE orders SET
		status = $1, accrual = $2, updated_at = now()
		WHERE order_id = '3498573'`
	_, err := op.db.Exec(ctx, query, status, orderID, accrual)
	if err != nil {
		return err
	}

	return nil
}

func (op *OrdersPostgres) Create(
	ctx context.Context, order models.OrderCreateDB, user models.DBUser,
) error {
	query := `INSERT INTO orders (order_id, accrual, status, user_login)
		values($1, $2, $3, $4)`

	_, err := op.db.Exec(ctx, query,
		order.OrderID, order.Accrual, order.Status, user.Login)
	if err != nil && pgutils.IsNotUnique(err) {
		oldOrder, err := op.GetOrderByOrderID(ctx, order.OrderID)
		if err != nil {
			return err
		}

		if oldOrder.UserLogin == user.Login {
			return errs.ErrOrdersRegisteredThisUser
		}

		return errs.ErrOrdersRegisteredOtherUser
	} else if err != nil {
		return err
	}

	return nil
}
