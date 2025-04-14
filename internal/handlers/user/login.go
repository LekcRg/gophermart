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
		logger.Log.Info(logContext + "user or password not found")
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

// Login godoc
// @Summary      Авторизация пользователя
// @Description  Авторизация по email и паролю, возвращает JWT токен при успешной аутентификации
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Данные для входа"
// @Success      200 {object} models.AuthResponse
// @Failure      400 {object} httputils.ErrorJSON "Неверные данные"
// @Failure      401 {object} httputils.ErrorJSON "Неверный email или пароль"
// @Failure      500 {object} httputils.ErrorJSON "Внутренняя ошибка сервера"
// @Router       /api/user/login [post]
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
		logger.Log.Error(logContext+"error while marshal json response",
			zap.Error(err))
		httputils.ErrInternalJSON(w)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
