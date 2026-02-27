package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/anditakaesar/uwa-go-fullstack/internal/server/transport"
	"github.com/anditakaesar/uwa-go-fullstack/internal/service"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xerror"
)

type UserApi struct {
	UserService service.IUserService
}

type UserApiDeps struct {
	UserService service.IUserService
}

func NewUserApi(dep UserApiDeps) *UserApi {
	return &UserApi{
		UserService: dep.UserService,
	}
}

func (h *UserApi) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var req CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return &xerror.ErrorBadRequest{Message: fmt.Sprintf("error while decoding request: %v", err)}
	}

	err = req.Validate()
	if err != nil {
		return err
	}

	user, err := h.UserService.CreateUser(r.Context(), domain.User{
		Username: strings.TrimSpace(req.Username),
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	transport.SendJSON(w, http.StatusCreated, UserDomainToResponse(user))
	return nil
}
