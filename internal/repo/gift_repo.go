package repo

import (
	"context"
	"fmt"

	"github.com/anditakaesar/uwa-go-fullstack/internal/common"
	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GiftRepository struct {
	db *pgxpool.Pool
}

func NewGiftRepository(db *pgxpool.Pool) *GiftRepository {
	return &GiftRepository{
		db: db,
	}
}

func (r *GiftRepository) Create(ctx context.Context, newGift domain.Gift) (*domain.Gift, error) {
	const query = `
		INSERT INTO gifts (title, description, stock, redeem_point, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, description, stock, redeem_point, image_url, created_at, updated_at, deleted_at;
	`

	var model domain.Gift

	err := r.db.QueryRow(ctx, query,
		newGift.Title,
		newGift.Description,
		newGift.Stock,
		newGift.RedeemPoint,
		newGift.ImageURL,
	).Scan(
		&model.ID,
		&model.Title,
		&model.Description,
		&model.Stock,
		&model.RedeemPoint,
		&model.ImageURL,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *GiftRepository) GetByID(ctx context.Context, id int64) (*domain.Gift, error) {
	const query = `
		SELECT id, title, description, stock, redeem_point, image_url, created_at, updated_at, deleted_at
		FROM gifts
		WHERE id = $1 AND deleted_at IS NULL;
	`

	var model domain.Gift

	err := r.db.QueryRow(ctx, query, id).Scan(
		&model.ID,
		&model.Title,
		&model.Description,
		&model.Stock,
		&model.RedeemPoint,
		&model.ImageURL,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *GiftRepository) FindAll(ctx context.Context, pagination common.Pagination, sort common.Sort) ([]domain.Gift, error) {
	const query = `
		SELECT id, title, description, stock, redeem_point, image_url, created_at, updated_at, deleted_at
		FROM gifts
		WHERE deleted_at IS NULL
		ORDER BY %s
		LIMIT $1 OFFSET $2;
	`
	rows, err := r.db.Query(ctx,
		fmt.Sprintf(query, sort.ToSQLSort()),
		pagination.Size, pagination.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("query gifts: %w", err)
	}
	defer rows.Close()

	gifts := []domain.Gift{}

	for rows.Next() {
		var g domain.Gift
		err := rows.Scan(
			&g.ID,
			&g.Title,
			&g.Description,
			&g.Stock,
			&g.RedeemPoint,
			&g.ImageURL,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		gifts = append(gifts, g)
	}

	return gifts, rows.Err()
}

func (r *GiftRepository) Count(ctx context.Context) (int64, error) {
	const query = `
		SELECT COUNT(*) FROM gifts WHERE deleted_at IS NULL;
	`
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
