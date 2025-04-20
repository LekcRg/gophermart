package orders

import (
	"context"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/validator"
)

var logContext = "OrderHandler"

type OrdersService interface {
	UploadOrder(ctx context.Context, order string) error
	GetOrders(ctx context.Context) ([]models.OrderDB, error)
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
