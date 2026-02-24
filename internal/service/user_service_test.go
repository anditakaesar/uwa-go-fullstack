package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/anditakaesar/uwa-go-fullstack/internal/mocks"
	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/stretchr/testify/assert"
)

type mockItems struct {
	ctx         context.Context
	userRepo    *mocks.MockIUserRepository
	passChecker *mocks.MockIPasswordChecker
	now         time.Time
}

func setupMocks() *mockItems {
	mockUserRepo := new(mocks.MockIUserRepository)
	mockPassChecker := new(mocks.MockIPasswordChecker)
	return &mockItems{
		ctx:         context.Background(),
		userRepo:    mockUserRepo,
		passChecker: mockPassChecker,
		now:         time.Now(),
	}
}

func TestNewUserService(test *testing.T) {
	test.Parallel()

	test.Run("success", func(t *testing.T) {
		m := setupMocks()

		got := service.NewUserService(m.userRepo, m.passChecker)
		assert.NotNil(t, got)
	})

}

func TestUserService_CreateUser(test *testing.T) {
	test.Parallel()

	userParam := domain.User{
		Username: "newusernonadmin",
		Password: "Some Pass",
	}

	test.Run("success", func(t *testing.T) {
		m := setupMocks()
		s := service.NewUserService(m.userRepo, m.passChecker)

		userReponse := domain.User{
			Base: domain.Base{
				ID:        1,
				CreatedAt: m.now,
			},
			Username: "John Doe",
			Role:     domain.RoleUser,
			Password: "somestring",
		}

		m.passChecker.On("HashPassword", userParam.Password).Return("somestring", nil).Once()
		updatedParam := userParam
		updatedParam.Password = "somestring"
		m.userRepo.On("CreateUser", m.ctx, updatedParam).Return(&userReponse, nil).Once()

		got, gotErr := s.CreateUser(m.ctx, userParam)
		assert.NoError(t, gotErr)

		assert.Equal(t, userReponse.Username, got.Username)
		m.passChecker.AssertExpectations(t)
		m.userRepo.AssertExpectations(t)
	})

	test.Run("error when hashing password", func(t *testing.T) {
		m := setupMocks()
		s := service.NewUserService(m.userRepo, m.passChecker)

		m.passChecker.On("HashPassword", userParam.Password).Return("", errors.New("error_HashPassword")).Once()

		got, gotErr := s.CreateUser(m.ctx, userParam)
		assert.Error(t, gotErr)
		assert.Nil(t, got)
		m.passChecker.AssertExpectations(t)
		m.userRepo.AssertExpectations(t)
	})

}

func TestUserService_CreateUserAdmin(test *testing.T) {
	test.Parallel()

	userParam := domain.User{
		Username: "newuseradmin",
		Password: "Some Pass",
	}

	test.Run("success", func(t *testing.T) {
		m := setupMocks()
		s := service.NewUserService(m.userRepo, m.passChecker)

		userReponse := domain.User{
			Base: domain.Base{
				ID:        1,
				CreatedAt: m.now,
			},
			Username: "John Doe",
			Role:     domain.RoleAdmin,
			Password: "somestring",
		}

		m.passChecker.On("HashPassword", userParam.Password).Return("somestring", nil).Once()
		updatedParam := userParam
		updatedParam.Password = "somestring"
		m.userRepo.On("CreateUserAdmin", m.ctx, updatedParam).Return(&userReponse, nil).Once()

		got, gotErr := s.CreateUserAdmin(m.ctx, userParam)
		assert.NoError(t, gotErr)

		assert.Equal(t, userReponse.Username, got.Username)
		m.passChecker.AssertExpectations(t)
		m.userRepo.AssertExpectations(t)
	})

	test.Run("error when hashing password", func(t *testing.T) {
		m := setupMocks()
		s := service.NewUserService(m.userRepo, m.passChecker)

		m.passChecker.On("HashPassword", userParam.Password).Return("", errors.New("error_HashPassword")).Once()

		got, gotErr := s.CreateUserAdmin(m.ctx, userParam)
		assert.Error(t, gotErr)
		assert.Nil(t, got)
		m.passChecker.AssertExpectations(t)
		m.userRepo.AssertExpectations(t)
	})

}

func TestUserService_AuthenticateUser(test *testing.T) {
	test.Parallel()

	test.Run("success", func(t *testing.T) {
		m := setupMocks()

		userResponse := domain.User{
			Username: "testuser",
			Password: "testpassword",
			Role:     domain.RoleAdmin,
			Base: domain.Base{
				ID: 1,
			},
		}

		m.userRepo.On("GetUser", m.ctx, userResponse.Username).Return(&userResponse, nil).Once()
		m.passChecker.On("CheckPassword", userResponse.Password, userResponse.Password).Return(true, nil).Once()

		s := service.NewUserService(m.userRepo, m.passChecker)
		got, gotErr := s.AuthenticateUser(m.ctx, userResponse.Username, userResponse.Password)

		assert.NoError(t, gotErr)
		assert.NotNil(t, got)
		m.userRepo.AssertExpectations(t)
		m.passChecker.AssertExpectations(t)
	})

	test.Run("password check failed", func(t *testing.T) {
		m := setupMocks()

		userResponse := domain.User{
			Username: "testuser",
			Password: "testpassword",
			Role:     domain.RoleAdmin,
			Base: domain.Base{
				ID: 1,
			},
		}

		m.userRepo.On("GetUser", m.ctx, userResponse.Username).Return(&userResponse, nil).Once()
		m.passChecker.On("CheckPassword", userResponse.Password, userResponse.Password).Return(false, nil).Once()

		s := service.NewUserService(m.userRepo, m.passChecker)
		got, gotErr := s.AuthenticateUser(m.ctx, userResponse.Username, userResponse.Password)

		assert.Error(t, gotErr)
		assert.Nil(t, got)
		m.userRepo.AssertExpectations(t)
		m.passChecker.AssertExpectations(t)
	})

	test.Run("error when get user", func(t *testing.T) {
		m := setupMocks()

		userResponse := domain.User{
			Username: "testuser",
			Password: "testpassword",
			Base: domain.Base{
				ID: 1,
			},
		}

		m.userRepo.On("GetUser", m.ctx, userResponse.Username).Return(nil, errors.New("error_GetUser")).Once()

		s := service.NewUserService(m.userRepo, m.passChecker)
		got, gotErr := s.AuthenticateUser(m.ctx, userResponse.Username, userResponse.Password)

		assert.Error(t, gotErr)
		assert.Nil(t, got)
		m.userRepo.AssertExpectations(t)
		m.passChecker.AssertExpectations(t)
	})
}

func TestUserService_GetUserByID(test *testing.T) {
	test.Parallel()

	test.Run("success", func(t *testing.T) {
		m := setupMocks()
		userResponse := domain.User{
			Username: "testuser",
			Password: "testpassword",
			Role:     domain.RoleAdmin,
			Base: domain.Base{
				ID: 1,
			},
		}

		m.userRepo.On("GetUserByID", m.ctx, int64(1)).Return(&userResponse, nil).Once()

		s := service.NewUserService(m.userRepo, m.passChecker)
		got, gotErr := s.GetUserByID(context.Background(), int64(1))

		assert.NoError(t, gotErr)
		assert.NotNil(t, got)
		assert.Equal(t, userResponse.Username, got.Username)
		m.userRepo.AssertExpectations(t)
	})

}
