package infra

import (
	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/repo"
	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Services struct {
	UserService   *service.UserService
	JWTService    *JWTService
	CookieService *CookieSvc
}

func NewInfra(pool *pgxpool.Pool) *Services {
	userRepo := repo.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo, &passwordUtil{})
	jwtSvc := NewJWTService(env.Values.JWTSecret)
	cookieService := NewCookieService(env.Values.IsDevelopment(), env.Values.CookieSecret)

	return &Services{
		UserService:   userSvc,
		JWTService:    jwtSvc,
		CookieService: cookieService,
	}
}
