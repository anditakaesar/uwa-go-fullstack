package infra

import (
	"context"

	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

type database struct {
	db *pgxpool.Pool
}

func NewDatabase() (*database, error) {
	pool, err := pgxpool.New(context.Background(), env.Values.DBUrl)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &database{
		db: pool,
	}, nil
}

func (d *database) Get() *pgxpool.Pool {
	return d.db
}

func (d *database) Close() {
	d.db.Close()
}
