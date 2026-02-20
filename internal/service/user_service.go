package service

import (
	"context"
	"fmt"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
)

type IUserService interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	AuthenticateUser(ctx context.Context, username string, password string) (*domain.User, error)
}

type UserService struct {
	userRepo    IUserRepository
	passChecker IPasswordChecker
}

func NewUserService(userRepo IUserRepository, passChecker IPasswordChecker) *UserService {
	return &UserService{
		userRepo:    userRepo,
		passChecker: passChecker,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	hash, err := s.passChecker.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hash
	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) AuthenticateUser(ctx context.Context, username string, password string) (*domain.User, error) {
	getUser, err := s.userRepo.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error while getting user: %v", err)
	}

	if !s.passChecker.CheckPassword(getUser.Password, password) {
		return nil, fmt.Errorf("wrong password attempt: %s", password)
	}

	return getUser, nil
}
