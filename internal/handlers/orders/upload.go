package orders

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/LekcRg/gophermart/internal/errs"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/validator"
	"go.uber.org/zap"
)

func (oh *OrdersHandler) checkUploadOrderError(
	w http.ResponseWriter, err error, lang string,
) {
	logContext := fmt.Sprintf("[%s/Login] ", logContext)
	var validErrs validator.ValidationErrors
	if err == errs.ErrOrdersRegisteredThisUser {
		logger.Log.Info(logContext, zap.Error(err))
		http.Error(w, err.Error(), http.StatusOK)
		return
	} else if err == errs.ErrOrdersRegisteredOtherUser {
		logger.Log.Info(logContext, zap.Error(err))
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if errors.As(err, &validErrs) &&
		len(validErrs) > 0 {
		trans := oh.validator.GetTrans(lang)
		tr, err := trans.T(validErrs[0].Tag())
		if err != nil {
			logger.Log.Error(logContext+"error getting translate",
				zap.Error(err))
			http.Error(w, "error", http.StatusUnprocessableEntity)
		}

		logger.Log.Info(logContext + "Bad order number")
		http.Error(w, tr, http.StatusUnprocessableEntity)
		return
	}

	logger.Log.Error(logContext+"internal error",
		zap.Error(err))
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// UploadOrder godoc
// @Summary      Загрузка номера заказа
// @Description  Загрузка заказов пользователя
// @Tags         Orders
// @Accept       text/plain
// @Produce      text/plain
// @Param        order-num body int true "Номер заказа"
// @Success      200 "Номер заказа уже был загружен этим пользователем"
// @Success      202 "Новый номер заказа принят в обработку"
// @Failure      400 "Неверный формат запроса"
// @Failure      401 "Пользователь не аутентифицирован"
// @Failure      409 "Номер заказа уже был загружен другим пользователем"
// @Failure      422 "Неверный формат номера заказа"
// @Failure      500 "Внутренняя ошибка сервера"
// @Router       /api/user/orders [post]
// @Security     BearerAuth
func (oh *OrdersHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		errText := "incorrect Content-Type"
		logger.Log.Info("[%s/UploadOrder] " + errText)
		http.Error(w, errText, http.StatusBadRequest)
		return
	}

	orderBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Log.Info("[%s/UploadOrder] "+"error while read body",
			zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	order := string(orderBytes)
	err = oh.service.UploadOrder(r.Context(), order)
	lang := r.Header.Get("Accept-Language")
	if err != nil {
		oh.checkUploadOrderError(w, err, lang)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Success"))
}
