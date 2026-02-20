package repo

import (
	"context"

	"github.com/anditakaesar/uwa-go-fullstack/internal/common"
	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db IDBExecutor
}

func NewUserRepository(db IDBExecutor) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetExecutor(ctx context.Context) IDBExecutor {
	tx, ok := ctx.Value(common.TxKey).(pgx.Tx)
	if ok {
		return tx
	}

	return r.db
}

func (r *UserRepository) CreateUser(ctx context.Context, newUser domain.User) (*domain.User, error) {
	const query = `
        INSERT INTO users (username, password)
        VALUES ($1, $2)
        RETURNING id, username, created_at, updated_at, deleted_at;
    `

	var model domain.User

	err := r.db.QueryRow(ctx, query, newUser.Username, newUser.Password).Scan(
		&model.ID,
		&model.Username,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *UserRepository) CreateUserAdmin(ctx context.Context, newUser domain.User) (*domain.User, error) {
	const query = `
        INSERT INTO users (username, password, role)
        VALUES ($1, $2, $3)
        RETURNING id, username, role, created_at, updated_at, deleted_at;
    `

	var model domain.User

	err := r.db.QueryRow(ctx, query, newUser.Username, newUser.Password, domain.RoleAdmin).Scan(
		&model.ID,
		&model.Username,
		&model.Role,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *UserRepository) GetUser(ctx context.Context, username string) (*domain.User, error) {
	const query = `
		SELECT id, username, password, role, created_at, updated_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
		AND username = $1
	`

	var model domain.User

	err := r.db.QueryRow(ctx, query, username).Scan(
		&model.ID,
		&model.Username,
		&model.Password,
		&model.Role,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	const query = `
		SELECT id, username, role, created_at, updated_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
		AND id = $1
	`

	var model domain.User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&model.ID,
		&model.Username,
		&model.Role,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
