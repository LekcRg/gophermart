package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/LekcRg/gophermart/internal/errs"
	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (uh *UserHandler) checkLoginErrors(w http.ResponseWriter, err error, lang string) {
	if errors.Is(err, pgx.ErrNoRows) ||
		errors.Is(err, errs.IncorrectPassword) ||
		errors.Is(err, errs.NotFoundUser) {
		logger.Log.Error(logContext+"user or password not found",
			zap.Error(err))
		httputils.ErrJSON(w, "user or password not found", http.StatusBadRequest)
		return
	} else if validErrs := uh.validator.GetValidTranslateErrs(err, lang); len(validErrs) > 0 {
		logger.Log.Info(logContext+"Bad validating login user",
			// maybe do pretty zap loggin validErrs
			zap.String("validation errors", fmt.Sprintf("%+v", validErrs)),
		)
		httputils.ErrMapJSON(w, validErrs, http.StatusBadRequest)

		return
	}

	logger.Log.Error(logContext+"error from service",
		zap.Error(err))
	httputils.ErrInternalJSON(w)
}

// JWT
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ0OTQ1OTgsIkxvZ2luIjoiJDJhJDEwJEFoNjZPL2RVZlJjQ0Zlc3pUNEN2ak9idmlMUGg2Lm5Gay5nQktQaDEuWE5GWkVidzRRRTZ1In0.46JyoQssLaqb1nnJYSLrTXQTwtXLX82SxWfzNZB4d-0
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	logContext := fmt.Sprintf("[%s/Login] ", logContext)
	var loginUser models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		logger.Log.Error(logContext + "invalid json")
		httputils.ErrJSON(w, "invalid JSON", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	token, err := uh.service.Login(r.Context(), loginUser)
	lang := r.Header.Get("Accept-Language")
	if err != nil {
		uh.checkLoginErrors(w, err, lang)
		return
	}

	res := models.AuthResponse{
		Token: token,
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		logger.Log.Error(logContext+"error while masrhal json response",
			zap.Error(err))
		httputils.ErrInternalJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resJson)
}
