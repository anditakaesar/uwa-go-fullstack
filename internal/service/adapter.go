package service

import (
	"context"

	"github.com/anditakaesar/uwa-go-fullstack/internal/common"
	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
)

type IGiftRepository interface {
	Create(ctx context.Context, newGift domain.Gift) (*domain.Gift, error)
	FindAll(ctx context.Context, pagination common.Pagination, sort common.Sort) ([]domain.Gift, error)
	Count(ctx context.Context) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Gift, error)
}

type IUserRepository interface {
	CreateUser(ctx context.Context, newUser domain.User) (*domain.User, error)
	GetUser(ctx context.Context, username string) (*domain.User, error)
}
