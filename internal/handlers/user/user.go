package user

import (
	"context"
	"net/http"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/validator"
)

type UserService interface {
	Register(ctx context.Context, user models.RegisterRequest) (string, error)
	Login(ctx context.Context, user models.LoginRequest) (string, error)
}

type UserHandler struct {
	service   UserService
	validator *validator.Validator
	config    config.Config
}

const logContext = "UserHandler"

func New(cfg config.Config, us UserService, validator *validator.Validator) *UserHandler {
	return &UserHandler{
		service:   us,
		validator: validator,
		config:    cfg,
	}
}

// IsAuth godoc
// @Summary      Проверка авторизации
// @Description  Проверяет авторизацию пользователя
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} httputils.MessageJSON
// @Failure      401 {object} httputils.ErrorJSON "Unauthorized"
// @Failure      500 {object} httputils.ErrorJSON "Internal server error"
// @Router       /api/user/is-auth [get]
// @Security     BearerAuth
func (us *UserHandler) IsAuth(w http.ResponseWriter, _ *http.Request) {
	httputils.SuccessJSON(w)
}
