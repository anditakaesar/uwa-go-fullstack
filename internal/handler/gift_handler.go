package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xlog"
	"github.com/go-chi/chi/v5"
)

const (
	gifts_endpoint string = "/gifts"
)

type GiftHandler struct {
	giftSvc service.IGiftService
}

func NewGiftHandler(giftSvc service.IGiftService) *GiftHandler {
	return &GiftHandler{
		giftSvc: giftSvc,
	}
}

func SetupGiftRoutes(router chi.Router, handler *GiftHandler) {
	endpoints := []EndpointWithMiddleware{
		{
			Endpoint: Endpoint{
				HttpMethod: http.MethodGet,
				Path:       gifts_endpoint,
				Handler:    handler.ListGifts,
			},
			Middlewares: []func(http.Handler) http.Handler{
				RequireAuth(),
			},
		},
		{
			Endpoint: Endpoint{
				HttpMethod: http.MethodGet,
				Path:       gifts_endpoint + "/{id}",
				Handler:    handler.GetGiftByID,
			},
			Middlewares: []func(http.Handler) http.Handler{
				RequireAuth(),
			},
		},
	}

	for _, e := range endpoints {
		if len(e.Middlewares) > 0 {
			router.With(e.Middlewares...).MethodFunc(e.HttpMethod, e.Path, e.Handler)
		} else {
			router.MethodFunc(e.HttpMethod, e.Path, e.Handler)
		}
	}
}

func (h *GiftHandler) ListGifts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pagination := ParsePagination(r)
	sort := ParseSort(r)

	gifts, total, err := h.giftSvc.ListGifts(ctx, pagination, sort)
	if err != nil {
		xlog.Logger.Error(fmt.Sprintf("error while fetching gifts: %v", err))
		JSONAPIErr(w, http.StatusInternalServerError, ErrObj{
			Title: "error while fething gifts",
		})
		return
	}

	resources := make([]GiftResource, 0, len(gifts))
	for _, g := range gifts {
		resources = append(resources, GiftDomainToResource(g))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Size)))

	response := GiftListResponse{
		Data:  resources,
		Links: BuildGiftListPaginationLinks(r, pagination, totalPages),
		Meta: GiftListMeta{
			Page: GiftListMetaPage{
				TotalPages: totalPages,
				Number:     pagination.Page,
				Size:       pagination.Size,
				TotalItems: total,
			},
		},
	}
	JSONAPI(w, http.StatusOK, response)
}

func (h *GiftHandler) GetGiftByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		JSONAPIErr(w, http.StatusBadRequest, ErrObj{
			Title: "missing id param",
			Source: ErrObjSource{
				Parameter: "id",
			},
		})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		JSONAPIErr(w, http.StatusBadRequest, ErrObj{
			Title: "invalid id param",
			Source: ErrObjSource{
				Parameter: "id",
			},
		})
		return
	}

	gift, err := h.giftSvc.GetByID(r.Context(), id)
	if err != nil {
		JSONAPIErr(w, http.StatusInternalServerError, ErrObj{
			Title: "failed to fetch gift by id",
			Source: ErrObjSource{
				Parameter: fmt.Sprintf("id=%d", id),
			},
		})
		return
	}

	JSONAPI(w, http.StatusOK, map[string]any{
		"data": GiftDomainToResource(*gift),
	})
}
