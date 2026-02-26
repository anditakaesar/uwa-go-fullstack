package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xerror"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

const (
	loginPage  string = "login.html"
	sessionKey string = "auth_session"
)

type MainHandler struct {
	UserService   service.IUserService
	JWTService    IJWTService
	CookieService ICookieService
	Render        func(context.Context, http.ResponseWriter, string, map[string]any)
}

type MainHandlerDeps struct {
	UserService   service.IUserService
	JWTService    IJWTService
	CookieService ICookieService
	WebRenderer   IWebRenderer
}

func NewMainHandler(dep MainHandlerDeps) *MainHandler {
	return &MainHandler{
		UserService:   dep.UserService,
		JWTService:    dep.JWTService,
		CookieService: dep.CookieService,
		Render:        dep.WebRenderer.Render2,
	}
}

func SetupMainRoutes(router chi.Router, handler *MainHandler) {
	endpoints := []Endpoint{
		{
			HttpMethod: http.MethodGet,
			Path:       "/",
			Handler:    handler.Index,
		},

		{
			HttpMethod: http.MethodGet,
			Path:       "/login",
			Handler:    handler.GetLogin,
		},

		{
			HttpMethod: http.MethodPost,
			Path:       "/login",
			Handler:    MakeHandler(handler.DoLogin),
		},
	}

	protectedEndpoints := []EndpointWithMiddleware{
		{
			Endpoint: Endpoint{
				HttpMethod: http.MethodGet,
				Path:       "/logout",
				Handler:    MakeHandler(handler.DoLogout),
			},
			Middlewares: []func(http.Handler) http.Handler{
				RequireAuth(),
				CSRFMiddleware(),
			},
		},
		{
			Endpoint: Endpoint{
				HttpMethod: http.MethodGet,
				Path:       "/upload",
				Handler:    handler.GetUploadPage,
			},
			Middlewares: []func(http.Handler) http.Handler{
				RequireAuth(),
				RequireRole([]domain.Role{domain.RoleAdmin}),
				CSRFMiddleware(),
			},
		},
		{
			Endpoint: Endpoint{
				HttpMethod: http.MethodPost,
				Path:       "/upload",
				Handler:    handler.PostUpload,
			},
			Middlewares: []func(http.Handler) http.Handler{
				RequireAuth(),
				CSRFMiddleware(),
			},
		},
	}

	router.Group(func(r chi.Router) {
		r.Use(CSRFMiddleware())
		for _, endpoint := range endpoints {
			r.MethodFunc(endpoint.HttpMethod, endpoint.Path, endpoint.Handler)
		}
	})

	for _, e := range protectedEndpoints {
		if len(e.Middlewares) > 0 {
			router.With(e.Middlewares...).MethodFunc(e.HttpMethod, e.Path, e.Handler)
		} else {
			router.MethodFunc(e.HttpMethod, e.Path, e.Handler)
		}
	}
}

func (h *MainHandler) Index(w http.ResponseWriter, r *http.Request) {
	session, err := h.CookieService.Get(r, sessionKey)
	if err != nil {
		SendError(w, http.StatusInternalServerError, ErrObj{
			Title:   "error when get auth_session",
			Message: err.Error(),
		})
		return
	}

	data := map[string]any{
		"Title": "Home Page",
		"Name":  "Index page html",
	}

	token, ok := session.Values["token"].(string)
	if ok {
		data["Token"] = token
	}

	h.Render(r.Context(), w, "index.html", data)
}

func (h *MainHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	h.Render(r.Context(), w, loginPage, map[string]any{
		"CSRF": csrf.Token(r),
	})
}

func (h *MainHandler) DoLogout(w http.ResponseWriter, r *http.Request) error {
	session, err := h.CookieService.Get(r, sessionKey)
	if err != nil {
		return &xerror.ErrorSession{Message: err.Error()}
	}

	delete(session.Values, "user_id")
	delete(session.Values, "username")
	delete(session.Values, "token")

	err = h.CookieService.Save(session, r, w)
	if err != nil {
		return &xerror.ErrorSession{Message: err.Error()}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (h *MainHandler) DoLogin(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return &xerror.ErrorBadRequest{Message: err.Error()}
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		h.Render(r.Context(), w, loginPage, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "username and password required",
		})
		return nil
	}

	user, err := h.UserService.AuthenticateUser(r.Context(), username, password)
	if err != nil {
		h.Render(r.Context(), w, loginPage, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "invalid credentials",
		})
		return nil
	}

	session, err := h.CookieService.Get(r, sessionKey)
	if err != nil {
		return &xerror.ErrorSession{Message: err.Error()}
	}
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username

	jwtToken, err := h.JWTService.IssueJWT(user.ID, []byte(env.Values.JWTSecret))
	if err != nil {
		return &xerror.ErrorToken{Message: err.Error()}
	}

	session.Values["token"] = jwtToken

	err = h.CookieService.Save(session, r, w)
	if err != nil {
		return &xerror.ErrorSession{Message: err.Error()}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (h *MainHandler) GetUploadPage(w http.ResponseWriter, r *http.Request) {
	h.Render(r.Context(), w, "upload.html", map[string]any{
		"CSRF": csrf.Token(r),
	})
}

func (h *MainHandler) PostUpload(w http.ResponseWriter, r *http.Request) {
	const uploadHTML = "upload.html"
	err := r.ParseMultipartForm(env.MAX_UPLOAD_SIZE)
	if err != nil {
		h.Render(r.Context(), w, uploadHTML, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "error when parsing file",
		})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.Render(r.Context(), w, uploadHTML, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "bad request",
		})
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		h.Render(r.Context(), w, uploadHTML, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "error while processing the file",
		})
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		h.Render(r.Context(), w, uploadHTML, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "error while processing the file",
		})
		return
	}

	if !env.UPLOAD_ALLOWED_TYPES[http.DetectContentType(buff)] {
		h.Render(r.Context(), w, uploadHTML, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "file type not allowed",
		})
		return
	}

	newName := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	dst, err := os.Create(env.Values.UploadDir + "/" + newName)
	if err != nil {
		h.Render(r.Context(), w, uploadHTML, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "error while performing save file request",
		})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		h.Render(r.Context(), w, uploadHTML, map[string]any{
			"CSRF":  csrf.Token(r),
			"Error": "error while performing copy file request",
		})
		return
	}

	h.Render(r.Context(), w, uploadHTML, map[string]any{
		"CSRF":     csrf.Token(r),
		"Uploaded": "uploads/" + newName,
	})
}
