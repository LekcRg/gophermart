package orders

import (
	"context"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/crypto"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/validator"
	"go.uber.org/zap"
)

type OrdersService struct {
	config    config.Config
	validator *validator.Validator
	db        repository.OrdersRepository
}

func New(
	db repository.OrdersRepository,
	validator *validator.Validator,
	config config.Config,
) *OrdersService {

	return &OrdersService{
		config:    config,
		validator: validator,
		db:        db,
	}
}

func (os *OrdersService) UploadOrder(
	ctx context.Context, order string,
) error {
	err := os.validator.ValidateVar(order, "required,luhn-order")

	if err != nil {
		return err
	}

	user, err := crypto.GetUserFromCtx(ctx)
	if err != nil {
		logger.Log.Error("error while getting user data from context")
		return err
	}

	// TODO: request
	// got status
	status := repository.OrderStatusNew

	err = os.db.Create(ctx, order, status, user)
	if err != nil {
		logger.Log.Error("order service db err",
			zap.Error(err))
	}

	return nil
}
