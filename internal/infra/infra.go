package infra

import (
	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/repo"
	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/anditakaesar/uwa-go-fullstack/internal/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Services struct {
	UserService   *service.UserService
	JWTService    *JWTService
	CookieService *CookieSvc
	WebRenderer   *web.Renderer
}

func NewInfra(pool *pgxpool.Pool) *Services {
	userRepo := repo.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo, NewPasswordHelper(env.Values.PassSecret))
	jwtSvc := NewJWTService(env.Values.JWTSecret)
	cookieService := NewCookieService(env.Values.IsDevelopment(), env.Values.CookieSecret)

	return &Services{
		UserService:   userSvc,
		JWTService:    jwtSvc,
		CookieService: cookieService,
		WebRenderer:   web.NewRenderer(),
	}
}
