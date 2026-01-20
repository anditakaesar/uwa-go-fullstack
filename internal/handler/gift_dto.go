package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anditakaesar/uwa-go-fullstack/internal/common"
	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
)

type GiftAttributes struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	RedeemPoint int       `json:"redeemPoint"`
	ImageURL    string    `json:"imageURL"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GiftResource struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Attributes GiftAttributes `json:"attributes"`
}

type GiftListResponse struct {
	Data  []GiftResource `json:"data"`
	Links GiftListLinks  `json:"links"`
	Meta  GiftListMeta   `json:"meta"`
}

type GiftListLinks struct {
	Self  *string `json:"self,omitempty"`
	First *string `json:"first,omitempty"`
	Last  *string `json:"last,omitempty"`
	Prev  *string `json:"prev,omitempty"`
	Next  *string `json:"next,omitempty"`
}

type GiftListMeta struct {
	Page GiftListMetaPage `json:"page"`
}

type GiftListMetaPage struct {
	Number     int   `json:"number,omitempty"`
	Size       int   `json:"size,omitempty"`
	TotalPages int   `json:"totalPages,omitempty"`
	TotalItems int64 `json:"totalItems,omitempty"`
}

func GiftDomainToResource(g domain.Gift) GiftResource {
	return GiftResource{
		Type: domain.GIFT_TYPE,
		ID:   strconv.FormatInt(g.ID, 10),
		Attributes: GiftAttributes{
			Title:       g.Title,
			Description: g.Description,
			Stock:       g.Stock,
			RedeemPoint: g.RedeemPoint,
			ImageURL:    fmt.Sprintf("%s/uploads/%s", env.Values.HostName, g.ImageURL),
			CreatedAt:   g.CreatedAt,
		},
	}
}

func BuildGiftListPaginationLinks(r *http.Request, pagination common.Pagination, totalPages int) GiftListLinks {
	urlTemplate := "%s?page[number]=%d&page[size]=%d"
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "" {
		urlTemplate += fmt.Sprintf("?sort=%s", sortParam)
	}

	resultLinks := GiftListLinks{}

	self := fmt.Sprintf(urlTemplate, gifts_endpoint, pagination.Page, pagination.Size)
	resultLinks.Self = &self

	first := fmt.Sprintf(urlTemplate, gifts_endpoint, 1, pagination.Size)
	resultLinks.First = &first

	last := fmt.Sprintf(urlTemplate, gifts_endpoint, totalPages, pagination.Size)
	resultLinks.Last = &last

	prevPage := pagination.Page - 1
	if prevPage > 0 {
		prev := fmt.Sprintf(urlTemplate, gifts_endpoint, prevPage, pagination.Size)
		resultLinks.Prev = &prev
	}

	nextPage := pagination.Page + 1
	if nextPage <= totalPages {
		next := fmt.Sprintf(urlTemplate, gifts_endpoint, nextPage, pagination.Size)
		resultLinks.Next = &next
	}

	return resultLinks
}

func ParsePagination(r *http.Request) common.Pagination {
	const (
		defaultPage     int = 1
		defaultPageSize int = 10
		maxPageSize     int = 50
	)

	q := r.URL.Query()
	page, err := strconv.Atoi(q.Get("page[number]"))
	if err != nil {
		page = defaultPage
	}

	size, err := strconv.Atoi(q.Get("page[size]"))
	if err != nil {
		size = defaultPageSize
	}

	if page < 1 {
		page = defaultPage
	}

	if size < 1 {
		size = defaultPageSize
	}

	if size > maxPageSize {
		size = maxPageSize
	}

	return common.Pagination{
		Page: page,
		Size: size,
	}
}

func ParseSort(r *http.Request) common.Sort {
	sortParam := r.URL.Query().Get("sort")

	// default sort
	if sortParam == "" {
		return common.Sort{
			Field:     "created_at",
			Direction: common.SORT_DESC,
		}
	}

	direction := common.SORT_ASC
	field := sortParam

	if strings.HasPrefix(sortParam, "-") {
		direction = common.SORT_DESC
		field = strings.TrimPrefix(sortParam, "-")
	}

	switch field {
	case "createdAt":
		field = "created_at"
	default:
		field = "created_at"
		direction = common.SORT_DESC
	}

	return common.Sort{
		Field:     field,
		Direction: direction,
	}
}
