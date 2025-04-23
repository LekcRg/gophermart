package withdraw

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/LekcRg/gophermart/internal/errs"
	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/service/withdraw"
	"github.com/LekcRg/gophermart/internal/validator"
	"go.uber.org/zap"
)

type WithdrawHandler struct {
	service   *withdraw.WithdrawService
	validator *validator.Validator
}

func New(withdrawService *withdraw.WithdrawService, validator *validator.Validator) *WithdrawHandler {
	return &WithdrawHandler{
		service:   withdrawService,
		validator: validator,
	}
}

func (wh *WithdrawHandler) CheckWithdrawErrors(
	w http.ResponseWriter, err error, lang string,
) {

	if validErrs := wh.validator.GetValidTranslateErrs(err, lang); len(validErrs) > 0 {
		logger.Log.Info("Bad validating withdraw",
			zap.String("validation errors", fmt.Sprintf("%+v", validErrs)),
		)
		if validErrs["order"] != "" {
			httputils.ErrMapJSON(w, validErrs, http.StatusUnprocessableEntity)
			return
		}

		httputils.ErrMapJSON(w, validErrs, http.StatusBadRequest)
		return
	} else if err == errs.ErrUserSmallBalance {
		logger.Log.Info("Small balance err")
		httputils.ErrJSON(w, "small balance", http.StatusPaymentRequired)
		return
	} else if err != nil {
		logger.Log.Info("Withdraw handler error", zap.Error(err))
		httputils.ErrJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// WithdrawCreate godoc
// @Summary      Запрос на списание средств
// @Description  Запрос на списание средств
// @Tags         Balance
// @Accept       json
// @Produce      json
// @Success      200 "Успешно"
// @Failure      400 {object} map[string]string "Error"
// @Failure      401 {object} httputils.ErrorJSON "Unauthorized"
// @Failure      402 {object} httputils.ErrorJSON "На счету недостаточно средств"
// @Failure      500 {object} httputils.ErrorJSON "Internal server error"
// @Router       /api/user/balance/withdraw [post]
// @Security     BearerAuth
func (wh *WithdrawHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		httputils.ErrJSON(w, "Incorrect Content-Type", http.StatusBadRequest)
		return
	}

	var withdraw models.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		httputils.ErrJSON(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := wh.service.Withdraw(r.Context(), withdraw)
	lang := r.Header.Get("Accept-Language")
	if err != nil {
		wh.CheckWithdrawErrors(w, err, lang)
		return
	}

	httputils.SuccessJSON(w)
}

// WithdrawList godoc
// @Summary      Получение информации о выводе средств
// @Description  Получение информации о выводе средств
// @Tags         Balance
// @Produce      json
// @Success      200 {object} models.WithdrawList "Успешно"
// @Success      204 "нет ни одного списания"
// @Failure      401 {object} httputils.ErrorJSON "Unauthorized"
// @Failure      500 {object} httputils.ErrorJSON "Internal server error"
// @Router       /api/user/withdrawals [get]
// @Security     BearerAuth
func (wh *WithdrawHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	list, err := wh.service.GetWithdrawals(r.Context())
	if err != nil {
		logger.Log.Error("Get withdrawals error", zap.Error(err))
		httputils.ErrInternalJSON(w)
		return
	}

	res, err := json.Marshal(list)
	if err != nil {
		logger.Log.Error("marshal list error",
			zap.Error(err))
		httputils.ErrInternalJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(list) <= 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(res)
}
