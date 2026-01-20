package infra

import (
	"net/http"

	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/gorilla/sessions"
)

type CookieSvc struct {
	cookieStore *sessions.CookieStore
}

func NewCookieService() *CookieSvc {
	cookieStore := sessions.NewCookieStore(
		[]byte(env.Values.CookieSecret),
	)

	cookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   !env.Values.IsDevelopment(), // true in prod
	}

	return &CookieSvc{
		cookieStore: cookieStore,
	}
}

func (s *CookieSvc) Get(r *http.Request, name string) (*sessions.Session, error) {
	return s.cookieStore.Get(r, name)
}
