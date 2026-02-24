package domain_test

import (
	"context"
	"testing"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserFromContext(test *testing.T) {
	test.Parallel()

	test.Run("success", func(t *testing.T) {
		ctx := context.Background()
		user := domain.User{
			Base: domain.Base{
				ID: 1,
			},
			Username: "stored user",
			Role:     domain.RoleAdmin,
		}
		newCtx := context.WithValue(
			ctx,
			domain.UserCtxKey,
			&user,
		)

		got, ok := domain.UserFromContext(newCtx)
		assert.Equal(t, true, ok)
		assert.Equal(t, int64(1), got.Base.ID)
		assert.Equal(t, domain.RoleAdmin, got.Role)
	})
}
