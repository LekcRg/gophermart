package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type Request struct {
	accrualAddr string
}

func New(accrualAddr string) *Request {

	return &Request{
		accrualAddr: accrualAddr,
	}
}

var ErrNotRegisteredOrder = fmt.Errorf("заказ не зарегистрирован в системе расчёта")

func (r *Request) GetAccrual(orderNum string) (models.AccrualRes, error) {
	client := resty.New()
	defer client.Close()

	res, err := client.R().
		Get(r.accrualAddr + "/api/orders/" + orderNum)
	if err != nil {
		if res.StatusCode() == http.StatusNoContent {
			logger.Log.Info("Get accrual no content")
			return models.AccrualRes{}, ErrNotRegisteredOrder
		}
		logger.Log.Error("Get accrual err",
			zap.Error(err))
		return models.AccrualRes{}, err
	}

	var resStruct models.AccrualRes
	err = json.Unmarshal(res.Bytes(), &resStruct)
	if err != nil {
		logger.Log.Error("err unmarshal", zap.Error(err))
		return models.AccrualRes{}, err
	}

	return resStruct, nil
}
