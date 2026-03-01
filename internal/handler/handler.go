package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anditakaesar/uwa-go-fullstack/internal/server/transport"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xerror"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xlog"
	"github.com/go-chi/chi/v5"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

type Endpoint struct {
	HttpMethod string
	Path       string
	Handler    func(w http.ResponseWriter, r *http.Request)
}

type EndpointWithMiddleware struct {
	Endpoint
	Middlewares []func(http.Handler) http.Handler
}

func MakeHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			xlog.Logger.Error(fmt.Sprintf("Server error [%s]: %v", r.URL.Path, err))
			transport.SendError(w, xerror.DefineStatusCode(err), transport.ErrObj{
				Title:   "server error",
				Message: err.Error(),
			})
		}
	}
}

func parseIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("invalid id param")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id param")
	}

	return id, nil
}
