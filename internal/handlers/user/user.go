package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/crypto"
	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/validator"
	"go.uber.org/zap"
)

type UserService interface {
	Register(ctx context.Context, user models.RegisterRequest) (string, error)
	Login(ctx context.Context, user models.LoginRequest) (string, error)
	Balance(ctx context.Context) (models.UserBalance, error)
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

// UserInfo godoc
// @Summary      Информация о пользователе
// @Description  Информация о пользователе возвращается id и логин
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} models.DBUser "User info"
// @Failure      401 {object} httputils.ErrorJSON "Unauthorized"
// @Failure      500 {object} httputils.ErrorJSON "Internal server error"
// @Router       /api/user/info [get]
// @Security     BearerAuth
func (uh *UserHandler) Info(w http.ResponseWriter, r *http.Request) {
	user, err := crypto.GetUserFromCtx(r.Context())
	if err != nil {
		logger.Log.Error("error while getting user data from context")
		httputils.ErrJSON(w, "Unauthorized", http.StatusUnauthorized)
	}

	res, err := json.Marshal(user)
	if err != nil {
		logger.Log.Error("error marsha JSON",
			zap.Error(err))
		httputils.ErrInternalJSON(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UserBalance godoc
// @Summary      Информация о балансе пользователя
// @Description  Информация о балансе пользователя
// @Tags         Balance
// @Produce      json
// @Success      200 {object} models.UserBalance "User info"
// @Failure      401 {object} httputils.ErrorJSON "Unauthorized"
// @Failure      500 {object} httputils.ErrorJSON "Internal server error"
// @Router       /api/user/balance [get]
// @Security     BearerAuth
func (uh *UserHandler) Balance(w http.ResponseWriter, r *http.Request) {
	userBalance, err := uh.service.Balance(r.Context())
	if err != nil {
		logger.Log.Error("error getting balance",
			zap.Error(err))
		httputils.ErrInternalJSON(w)
	}

	res, err := json.Marshal(userBalance)
	if err != nil {
		logger.Log.Error("error marshal user balance",
			zap.Error(err))
		httputils.ErrInternalJSON(w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
