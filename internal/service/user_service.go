package service

import (
	"context"
	"fmt"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	AuthenticateUser(ctx context.Context, username string, password string) (*domain.User, error)
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	user.Password = s.hashPassword(user.Password)
	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) AuthenticateUser(ctx context.Context, username string, password string) (*domain.User, error) {
	getUser, err := s.userRepo.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error while getting user: %v", err)
	}

	if !s.checkPassword(getUser.Password, password) {
		return nil, fmt.Errorf("wrong password attempt: %s", password)
	}

	return getUser, nil
}

func (s *UserService) hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func (s *UserService) checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
