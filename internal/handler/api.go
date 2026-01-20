package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type ErrObjSource struct {
	Parameter string `json:"parameter,omitempty"`
	Pointer   string `json:"pointer,omitempty"`
	Header    string `json:"header,omitempty"`
}

type ErrObj struct {
	Status string       `json:"status"`
	Title  string       `json:"title"`
	Source ErrObjSource `json:"source"`
}

type ErrObjResponse struct {
	ErrObjs []ErrObj `json:"errors"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	err := enc.Encode(v)
	if err != nil {
		http.Error(w, "json encode error", http.StatusInternalServerError)
	}
}

func JSONAPI(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	err := enc.Encode(v)

	if err != nil {
		http.Error(w, "json:api encode error", http.StatusInternalServerError)
	}
}

func JSONAPIErr(w http.ResponseWriter, status int, errObj ErrObj) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	errObj.Status = fmt.Sprint(status)

	err := enc.Encode(ErrObjResponse{
		ErrObjs: []ErrObj{errObj},
	})

	if err != nil {
		http.Error(w, "json:api encode error", http.StatusInternalServerError)
	}
}
