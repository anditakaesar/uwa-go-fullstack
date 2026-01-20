package repo

import (
	"context"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
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

func (r *UserRepository) GetUser(ctx context.Context, username string) (*domain.User, error) {
	const query = `
		SELECT id, username, password, created_at, updated_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
		AND username = $1
	`

	var model domain.User

	err := r.db.QueryRow(ctx, query, username).Scan(
		&model.ID,
		&model.Username,
		&model.Password,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
