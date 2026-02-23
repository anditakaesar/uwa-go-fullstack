package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anditakaesar/uwa-go-fullstack/internal/handler"
	"github.com/anditakaesar/uwa-go-fullstack/internal/mocks"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockItems struct {
	ctx       context.Context
	cookieSvc *mocks.MockICookieService
}

func setupMocks() *mockItems {
	return &mockItems{
		ctx:       context.Background(),
		cookieSvc: new(mocks.MockICookieService),
	}
}

func TestMainHandler_Index(test *testing.T) {
	test.Parallel()

	test.Run("success", func(t *testing.T) {
		m := setupMocks()

		h := &handler.MainHandler{
			CookieService: m.cookieSvc,
			Render: func(ctx context.Context, w http.ResponseWriter, s string, m map[string]any) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "rendered")
			},
		}

		sessionMock := &sessions.Session{
			Values: map[any]any{
				"user_id":  "1",
				"username": "user1",
			},
		}

		m.cookieSvc.On("Get", mock.Anything, "auth_session").Return(sessionMock, nil).Once()

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		h.Index(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "rendered")
		m.cookieSvc.AssertExpectations(t)
	})

	test.Run("success and token found", func(t *testing.T) {
		m := setupMocks()

		h := &handler.MainHandler{
			CookieService: m.cookieSvc,
			Render: func(ctx context.Context, w http.ResponseWriter, s string, m map[string]any) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "rendered")
			},
		}

		sessionMock := &sessions.Session{
			Values: map[any]any{
				"user_id":  "1",
				"username": "user1",
				"token":    "some-token-here",
			},
		}

		m.cookieSvc.On("Get", mock.Anything, "auth_session").Return(sessionMock, nil).Once()

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		h.Index(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "rendered")
		m.cookieSvc.AssertExpectations(t)
	})

	test.Run("error when get session", func(t *testing.T) {
		m := setupMocks()

		h := &handler.MainHandler{
			CookieService: m.cookieSvc,
			Render: func(ctx context.Context, w http.ResponseWriter, s string, m map[string]any) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "rendered")
			},
		}

		m.cookieSvc.On("Get", mock.Anything, "auth_session").Return(nil, errors.New("error_GetSession")).Once()

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		h.Index(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		m.cookieSvc.AssertExpectations(t)
	})
}
