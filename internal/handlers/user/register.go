package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/pgutils"
	"github.com/LekcRg/gophermart/internal/translations"
	"go.uber.org/zap"
)

func (uh *UserHandler) checkRegisterErrors(
	w http.ResponseWriter, err error, lang string,
) {
	logContext := fmt.Sprintf("[%s/checkRegisterErrors] ", logContext)
	if pgutils.IsNotUnique(err) {
		trErr, err := translations.GetErr(translations.ErrAlreadyExists, "login", lang)
		if err != nil {
			logger.Log.Error("Error getting translate err",
				zap.Error(err))
			trErr = "already exists"
		}
		logger.Log.Info(logContext + "No unique user")
		httputils.ErrMapJSON(w, map[string]string{
			"login": trErr,
		}, http.StatusConflict)

		return
	} else if validErrs := uh.validator.GetValidTranslateErrs(err, lang); len(validErrs) > 0 {
		logger.Log.Info(logContext+"Bad validating register user",
			// maybe do pretty zal loggin validErrs
			zap.String("validation errors", fmt.Sprintf("%+v", validErrs)),
		)
		httputils.ErrMapJSON(w, validErrs, http.StatusBadRequest)

		return
	}

	logger.Log.Error(logContext+"Error register user", zap.Error(err))
	httputils.ErrInternalJSON(w)
}

func (uh *UserHandler) Register(
	w http.ResponseWriter, r *http.Request,
) {
	logContext := fmt.Sprintf("[%s/Register] ", logContext)
	var authUser models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&authUser); err != nil {
		logger.Log.Error(logContext+"invalid JSON",
			zap.Error(err))
		httputils.ErrJSON(w, "invalid JSON", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	token, err := uh.service.Register(r.Context(), authUser)
	lang := r.Header.Get("Accept-Language")
	if err != nil {
		uh.checkRegisterErrors(w, err, lang)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	msgStruct := models.AuthResponse{
		Token: token,
	}

	msg, err := json.Marshal(msgStruct)
	if err != nil {
		logger.Log.Error(logContext+"Error marshaling auth response", zap.Error(err))
		httputils.ErrInternalJSON(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}
