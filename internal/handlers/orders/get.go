package orders

import (
	"encoding/json"
	"net/http"

	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/logger"
	"go.uber.org/zap"
)

// UploadOrder godoc
// @Summary      Загрузка заказов пользователя
// @Description  Загрузка заказов пользователя
// @Tags         Orders
// @Produce      json
// @Success      200 "Номер заказа уже был загружен этим пользователем"
// @Success      204 "нет данных для ответа"
// @Failure      500 "Внутренняя ошибка сервера"
// @Router       /api/user/orders [get]
// @Security     BearerAuth
func (oh *OrdersHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := oh.service.GetOrders(r.Context())
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		logger.Log.Error("get orders err",
			zap.Error(err))
		httputils.ErrJSON(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	res, err := json.Marshal(orders)
	if err != nil {
		logger.Log.Error("get orders err",
			zap.Error(err))
		httputils.ErrJSON(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(res)
}
