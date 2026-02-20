package handler

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

type ICookieService interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
}

type IJWTService interface {
	Verify(token string) (domain.UserClaims, error)
	IssueJWT(userID int64, secret []byte) (string, error)
}

type Middleware func(http.Handler) http.Handler

func CSRFMiddleware() Middleware {
	secure := !env.Values.IsDevelopment()

	opts := []csrf.Option{
		csrf.FieldName(env.CSRF_TOKEN_FIELD_NAME),
		csrf.Secure(secure),
	}

	if !secure {
		opts = append(opts,
			csrf.TrustedOrigins([]string{
				"localhost" + env.Values.Port,
			}),
		)
	}

	return csrf.Protect(
		[]byte(env.Values.CSRFSecret),
		opts...,
	)
}

func ResolveAuth(
	cookieStore ICookieService,
	userService service.IUserService,
	jwtService IJWTService,
) Middleware {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := cookieStore.Get(r, "auth_session")
			if err == nil {
				uid, ok := session.Values["user_id"].(int64)
				if ok {
					ctx := context.WithValue(
						r.Context(),
						domain.IdentityKey,
						domain.Identity{UserID: uid, Method: "session"},
					)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tokenStr := strings.TrimPrefix(auth, "Bearer ")
				claims, err := jwtService.Verify(tokenStr)
				if err == nil {
					ctx := context.WithValue(
						r.Context(),
						domain.IdentityKey,
						domain.Identity{UserID: claims.UserID, Method: "jwt"},
					)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ResolveUser(
	userService service.IUserService,
) Middleware {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity, ok := r.Context().Value(domain.IdentityKey).(domain.Identity)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			user, _ := userService.GetUserByID(r.Context(), identity.UserID)
			if user != nil {
				ctx := context.WithValue(
					r.Context(),
					domain.UserCtxKey,
					user,
				)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireRole(roles []domain.Role) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := domain.UserFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if user != nil && !slices.Contains(roles, user.Role) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := r.Context().Value(domain.IdentityKey).(domain.Identity)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
