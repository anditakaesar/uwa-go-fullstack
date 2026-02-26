package handler

import (
	"net/http"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/go-chi/chi/v5"
)

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
