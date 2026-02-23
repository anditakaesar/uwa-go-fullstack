package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anditakaesar/uwa-go-fullstack/internal/xlog"
)

type TypeHandlerFunc func(w http.ResponseWriter, r *http.Request)

type Endpoint struct {
	HttpMethod string
	Path       string
	Handler    func(w http.ResponseWriter, r *http.Request)
}

type EndpointWithMiddleware struct {
	Endpoint
	Middlewares []func(http.Handler) http.Handler
}

type Response struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

// ErrorResponse
type ErrObj struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error ErrObj `json:"error"`
}

func SendJSON(w http.ResponseWriter, status int, unwrappedData any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	response := Response{Data: unwrappedData}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		xlog.Logger.Error(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func SendError(w http.ResponseWriter, status int, errObj ErrObj) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: errObj,
	}
	json.NewEncoder(w).Encode(resp)
}
