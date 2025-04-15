package request

import (
	"fmt"

	"resty.dev/v3"
)

type Order struct {
	Order string `json:"order"`
}

func Post() {
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
		SetBody(Order{
			Order: "3498573",
		}).
		Post("http://localhost:8888/api/orders")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

}
