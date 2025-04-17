package request

import (
	"fmt"
	"io"

	"github.com/LekcRg/gophermart/internal/logger"
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

func (r *Request) GetAccrual(orderNum string) {
	// HTTP, REST Client
	client := resty.New()
	defer client.Close()

	// res, err := client.R().
	// 	EnableTrace().
	// 	Get("https://httpbin.org/get")
	// fmt.Println(err)
	// fmt.Println(res)
	// fmt.Println(res.Request.TraceInfo())

	res, err := client.R().
		Get(r.accrualAddr + "/api/orders/" + orderNum)
	if err != nil {
		logger.Log.Error("Get accrual err",
			zap.Error(err))
		return
	}
	fmt.Println(res)
	fmt.Println(res.Status())
	// var body map[string]string
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Log.Error("read response accrual body err",
			zap.Error(err))
	}
	defer res.Body.Close()
	fmt.Println(string(body))
}
