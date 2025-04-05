package user

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/LekcRg/gophermart/internal/common/httputils"
	"github.com/LekcRg/gophermart/internal/common/pgutils"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"go.uber.org/zap"
)

type UserService interface {
	Register(ctx context.Context, user models.AuthRequest) error
}

type UserHandler struct {
	service UserService
}

func New(us UserService) *UserHandler {
	return &UserHandler{
		service: us,
	}
}

func (uh *UserHandler) Register(
	w http.ResponseWriter, r *http.Request,
) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type", http.StatusBadRequest)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("User Register handler, read body err",
			zap.Error(err))
		httputils.ErrInternalJSON(w)
	}
	defer r.Body.Close()

	var authUser models.AuthRequest
	err = json.Unmarshal(body, &authUser)
	if err != nil {
		logger.Log.Error("Error register user, unmarshal json",
			zap.Error(err))
		httputils.ErrInternalJSON(w)
		return
	}

	err = uh.service.Register(r.Context(), authUser)
	if err != nil {
		if pgutils.IsNotUnique(err) {
			httputils.ErrJSON(w, "Login already exists", http.StatusConflict)
		} else {
			logger.Log.Error("Error register user, db request",
				zap.Error(err))
			httputils.ErrInternalJSON(w)
		}
		return
	}

	httputils.SuccessJSON(w)
}
