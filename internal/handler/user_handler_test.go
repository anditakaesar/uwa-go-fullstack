package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/anditakaesar/uwa-go-fullstack/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestUserApi_CreateUser(test *testing.T) {
	test.Parallel()

	test.Run("success create user", func(t *testing.T) {
		m, d := setupMocks()

		h := handler.NewUserApi(handler.UserApiDeps{
			UserService: d.UserService,
		})

		m.userSvc.On("CreateUser", m.anything, domain.User{
			Username: "newuser",
			Password: "newpassword",
		}).Return(&domain.User{
			Base: domain.Base{
				ID: 1,
			},
			Username: "newuser",
			Role:     domain.RoleUser,
		}, nil).Once()

		userReq := `{"username":"newuser","password":"newpassword"}`

		req, err := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(userReq))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		gotErr := h.CreateUser(rr, req)

		assert.NoError(t, gotErr)
		assert.Equal(t, http.StatusCreated, rr.Code)
		m.userSvc.AssertExpectations(t)
	})

	test.Run("error when create user", func(t *testing.T) {
		m, d := setupMocks()

		h := handler.NewUserApi(handler.UserApiDeps{
			UserService: d.UserService,
		})

		m.userSvc.On("CreateUser", m.anything, domain.User{
			Username: "newuser",
			Password: "newpassword",
		}).Return(nil, errors.New("error_CreateUser")).Once()

		userReq := `{"username":"newuser","password":"newpassword"}`

		req, err := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(userReq))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		gotErr := h.CreateUser(rr, req)
		assert.Error(t, gotErr)
		m.userSvc.AssertExpectations(t)
	})

	test.Run("error when validate request", func(t *testing.T) {
		m, d := setupMocks()

		h := handler.NewUserApi(handler.UserApiDeps{
			UserService: d.UserService,
		})

		userReq := `{"username":"","password":"password"}`

		req, err := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(userReq))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		gotErr := h.CreateUser(rr, req)
		assert.Error(t, gotErr)
		m.userSvc.AssertExpectations(t)
	})

	test.Run("error when validate request", func(t *testing.T) {
		m, d := setupMocks()

		h := handler.NewUserApi(handler.UserApiDeps{
			UserService: d.UserService,
		})

		userReq := `{"username":"username","password":""}`

		req, err := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(userReq))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		gotErr := h.CreateUser(rr, req)
		assert.Error(t, gotErr)
		m.userSvc.AssertExpectations(t)
	})

	test.Run("error when decoding request", func(t *testing.T) {
		m, d := setupMocks()

		h := handler.NewUserApi(handler.UserApiDeps{
			UserService: d.UserService,
		})

		userReq := `{"username":"","password":"x}`

		req, err := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(userReq))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		gotErr := h.CreateUser(rr, req)
		assert.Error(t, gotErr)
		m.userSvc.AssertExpectations(t)
	})
}
