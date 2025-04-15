package postgres

import (
	"context"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/repository/postgres/orders"
	"github.com/LekcRg/gophermart/internal/repository/postgres/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Postgres struct {
	db    *pgxpool.Pool
	repos *repository.Repository
}

func New(ctx context.Context, cfg config.Config) *Postgres {
	conn, err := pgxpool.New(ctx, cfg.DBURI)
	if err != nil {
		logger.Log.Error("pgxpool conn err",
			zap.Error(err))
	}

	return &Postgres{
		db: conn,
		repos: &repository.Repository{
			User:   user.New(ctx, conn),
			Orders: orders.New(ctx, conn),
		},
	}
}

func (p *Postgres) GetRepositories() *repository.Repository {
	return p.repos
}

func (p *Postgres) Close() {
	p.db.Close()
}
