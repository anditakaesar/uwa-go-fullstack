package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/anditakaesar/uwa-go-fullstack/internal/common"
	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xlog"
	"github.com/jackc/pgx/v5"
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

type unitOfWork struct {
	db *pgxpool.Pool
}

func (u *unitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			xlog.Logger.Error(fmt.Sprintf("rollback err: %v", rollbackErr))
		}
	}()

	txCtx := context.WithValue(ctx, common.TxKey, tx)

	err = fn(txCtx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func NewUnitOfWork(db *pgxpool.Pool) *unitOfWork {
	return &unitOfWork{db: db}
}
