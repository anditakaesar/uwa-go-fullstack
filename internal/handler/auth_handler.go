package handler

import (
	"net/http"

	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/anditakaesar/uwa-go-fullstack/internal/web"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

type AuthHandler struct {
	UserService   service.IUserService
	JWTService    IJWTService
	CookieService ICookieService
	Render        func(http.ResponseWriter, string, any)
}

type MainHandlerDeps struct {
	UserService   service.IUserService
	JWTService    IJWTService
	CookieService ICookieService
}

func NewMainHandler(dep MainHandlerDeps) *AuthHandler {
	renderer := web.NewRenderer()
	return &AuthHandler{
		UserService:   dep.UserService,
		JWTService:    dep.JWTService,
		CookieService: dep.CookieService,
		Render:        renderer.Render,
	}
}

func SetupMainRoutes(router chi.Router, handler *AuthHandler) {
	endpoints := []Endpoint{
		{
			HttpMethod: http.MethodGet,
			Path:       "/",
			Handler:    handler.GetLogin,
		},
		{
			HttpMethod: http.MethodPost,
			Path:       "/login",
			Handler:    handler.PostLogin,
		},
	}

	router.Group(func(r chi.Router) {
		r.Use(CSRFMiddleware())
		for _, endpoint := range endpoints {
			r.MethodFunc(endpoint.HttpMethod, endpoint.Path, endpoint.Handler)
		}
	})
}

type LoginView struct {
	CSRF  string
	Error string
}

func (h *AuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	h.Render(w, "login.html", LoginView{
		CSRF: csrf.Token(r),
	})
}

func (h *AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		h.Render(w, "login.html", LoginView{
			CSRF:  csrf.Token(r),
			Error: "username and password required",
		})
		return
	}

	user, err := h.UserService.AuthenticateUser(r.Context(), username, password)
	if err != nil {
		h.Render(w, "login.html", LoginView{
			CSRF:  csrf.Token(r),
			Error: "invalid credentials",
		})
		return
	}

	session, _ := h.CookieService.Get(r, "auth_session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username

	if err := session.Save(r, w); err != nil {
		JSONAPIErr(w, http.StatusInternalServerError, ErrObj{
			Title: "session error",
		})
		return
	}

	jwtToken, err := h.JWTService.IssueJWT(user.ID, []byte(env.Values.JWTSecret))
	if err := session.Save(r, w); err != nil {
		JSONAPIErr(w, http.StatusInternalServerError, ErrObj{
			Title: "issuing token error",
		})
		return
	}

	JSON(w, 200, map[string]any{
		"data": map[string]string{
			"token": jwtToken,
		},
	})
}
