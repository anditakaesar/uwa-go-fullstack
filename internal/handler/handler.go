package handler

import (
	"fmt"
	"net/http"

	"github.com/anditakaesar/uwa-go-fullstack/internal/server/transport"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xerror"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xlog"
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
