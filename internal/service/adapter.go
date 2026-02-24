package service

import (
	"context"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, newUser domain.User) (*domain.User, error)
	CreateUserAdmin(ctx context.Context, newUser domain.User) (*domain.User, error)
	GetUser(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}

type IUnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type IPasswordChecker interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, hash string) (bool, error)
}
