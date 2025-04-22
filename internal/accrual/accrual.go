package accrual

import (
	"encoding/json"
	"net/http"

	"github.com/LekcRg/gophermart/internal/errs"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type Accrual struct {
	client      *resty.Client
	accrualAddr string
}

func New(accrualAddr string) *Accrual {

	return &Accrual{
		client:      resty.New(),
		accrualAddr: accrualAddr,
	}
}

func (a *Accrual) GetAccrual(orderNum string) (models.AccrualRes, error) {
	res, err := a.client.R().
		Get(a.accrualAddr + "/api/orders/" + orderNum)
	if err != nil {
		if res.StatusCode() == http.StatusNoContent {
			logger.Log.Info("Get accrual no content")
			return models.AccrualRes{}, errs.ErrAccrualReqNotRegisteredOrder
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

func (a *Accrual) Close() {
	a.client.Close()
}
