package infra

import (
	"github.com/anditakaesar/uwa-go-fullstack/internal/repo"
	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Services struct {
	UserService   *service.UserService
	GiftService   *service.GiftService
	JWTService    *JWTService
	CookieService *CookieSvc
}

func NewInfra(pool *pgxpool.Pool) *Services {
	userRepo := repo.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo)
	giftRepo := repo.NewGiftRepository(pool)
	giftSvc := service.NewGiftService(giftRepo)
	jwtSvc := NewJWTService()
	cookieService := NewCookieService()

	return &Services{
		UserService:   userSvc,
		GiftService:   giftSvc,
		JWTService:    jwtSvc,
		CookieService: cookieService,
	}
}
