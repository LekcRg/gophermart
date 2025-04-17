package orders

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/validator"
	"go.uber.org/zap"
)

type OrdersService interface {
	UploadOrder(ctx context.Context, order string) error
}

type OrdersHandler struct {
	config    config.Config
	validator *validator.Validator
	service   OrdersService
}

func New(
	config config.Config, service OrdersService,
	validator *validator.Validator,
) *OrdersHandler {

	return &OrdersHandler{
		config:    config,
		service:   service,
		validator: validator,
	}
}

func (oh *OrdersHandler) checkUploadOrderError(
	w http.ResponseWriter, err error, lang string,
) {
	var validErrs validator.ValidationErrors
	if err == repository.ErrOrdersRegisteredThisUser {
		http.Error(w, err.Error(), http.StatusOK)
	} else if err == repository.ErrOrdersRegisteredOtherUser {
		http.Error(w, err.Error(), http.StatusConflict)
	} else if errors.As(err, &validErrs) &&
		len(validErrs) > 0 {
		trans := oh.validator.GetTrans(lang)
		tr, err := trans.T(validErrs[0].Tag())
		if err != nil {
			logger.Log.Error("[%s/checkUploadOrderError] "+"error getting translate",
				zap.Error(err))
			http.Error(w, "error", http.StatusUnprocessableEntity)
		}

		logger.Log.Info("[%s/checkUploadOrderError] " + "Bad order number")
		http.Error(w, tr, http.StatusUnprocessableEntity)
		return
	}
}

// UploadOrder godoc
// @Summary      Загрузка номера заказа
// @Description  Доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
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
