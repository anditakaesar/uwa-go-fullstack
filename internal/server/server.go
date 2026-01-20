package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/handler"
	"github.com/anditakaesar/uwa-go-fullstack/internal/infra"
	"github.com/anditakaesar/uwa-go-fullstack/internal/web"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xlog"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IDatabase interface {
	Get() *pgxpool.Pool
	Close()
}

type ServerDependency struct {
	DB IDatabase
}

func SetupServer(dep *ServerDependency) *chi.Mux {
	router := chi.NewRouter()
	infraSvc := infra.NewInfra(dep.DB.Get())

	// static files
	sub, err := fs.Sub(web.PublicFS, "public")
	if err != nil {
		xlog.Logger.Error(fmt.Sprintf("static file sub failed: %v", err))
		os.Exit(1)
	}
	router.Handle(
		"/static/*",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.FS(sub)),
		),
	)

	router.Handle(
		"/uploads/*",
		http.StripPrefix(
			"/uploads/",
			http.FileServer(http.Dir(env.Values.UploadDir)),
		),
	)

	// handlers and routes
	mainHandler := handler.NewMainHandler(handler.MainHandlerDeps{
		UserService:   infraSvc.UserService,
		JWTService:    infraSvc.JWTService,
		CookieService: infraSvc.CookieService,
	})
	giftHandler := handler.NewGiftHandler(infraSvc.GiftService)

	router.Group(func(r chi.Router) {
		// middlewares
		r.Use(handler.ResolveAuth(
			infraSvc.CookieService,
			infraSvc.UserService,
			infraSvc.JWTService,
		))

		handler.SetupMainRoutes(r, mainHandler)
		handler.SetupGiftRoutes(r, giftHandler)
	})

	return router
}
