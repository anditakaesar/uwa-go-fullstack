package service

import (
	"context"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, newUser domain.User) (*domain.User, error)
	GetUser(ctx context.Context, username string) (*domain.User, error)
}
