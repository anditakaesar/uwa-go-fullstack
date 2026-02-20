package infra

import (
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
	jwtSvc := NewJWTService()
	cookieService := NewCookieService()

	return &Services{
		UserService:   userSvc,
		JWTService:    jwtSvc,
		CookieService: cookieService,
	}
}
