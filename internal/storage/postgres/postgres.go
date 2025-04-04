package postgres

import (
	"context"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Postgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) *Postgres {
	conn, err := pgxpool.New(ctx, cfg.DBURI)
	if err != nil {
		logger.Log.Error("pgxpool conn err",
			zap.Error(err))
	}

	req := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		LOGIN VARCHAR(20) UNIQUE NOT NULL,
		salt VARCHAR(10) NOT NULL,
		passhash varchar(32) NOT NULL
	)`
	conn.Exec(ctx, req)

	return &Postgres{
		db: conn,
	}
}

func (p *Postgres) Close() {
	p.db.Close()
}
