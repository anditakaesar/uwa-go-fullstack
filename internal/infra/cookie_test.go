package infra

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCookieService(test *testing.T) {
	test.Parallel()

	test.Run("success", func(t *testing.T) {
		svc := NewCookieService(true, "test-secret")

		assert.Equal(t, false, svc.cookieStore.Options.Secure)
	})
}

func TestCookieSvc_Get(test *testing.T) {
	test.Run("success", func(t *testing.T) {
		svc := NewCookieService(true, "very-secret-key")
		req := httptest.NewRequest("GET", "http://localhost", nil)

		session, err := svc.Get(req, "session-name")
		if err != nil {
			t.Fatalf("Failed to get session: %v", err)
		}

		if session.Name() != "session-name" {
			t.Errorf("Expected session name 'session-name', got %s", session.Name())
		}
	})
}
