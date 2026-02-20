package domain

import (
	"context"

	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
)

type User struct {
	Base
	Username string
	Password string
	Role     Role
}

type ctxKeyUser string

const UserCtxKey ctxKeyUser = env.USER_CTX_KEY

func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserCtxKey).(*User)
	return user, ok
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "default"
)
