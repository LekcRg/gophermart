package orders

import (
	"context"
	"time"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/jwt"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/request"
	"github.com/LekcRg/gophermart/internal/validator"
	"go.uber.org/zap"
)

type OrdersService struct {
	config    config.Config
	validator *validator.Validator
	db        repository.OrdersRepository
	request   *request.Request
	userDB    repository.UserRepository
}

func New(
	db repository.OrdersRepository,
	validator *validator.Validator,
	config config.Config, req *request.Request,
	userDB repository.UserRepository,
) *OrdersService {

	return &OrdersService{
		config:    config,
		validator: validator,
		request:   req,
		db:        db,
		userDB:    userDB,
	}
}

func (os *OrdersService) GetOrders(
	ctx context.Context,
) ([]models.OrderDB, error) {
	user, err := jwt.GetUserFromCtx(ctx)
	if err != nil {
		logger.Log.Error("error while getting user data from context")
		return []models.OrderDB{}, err
	}
	return os.db.GetOrdersByUserLogin(ctx, user.Login)
}

// context ???
func (os *OrdersService) StartMonitoringStatus(
	order models.OrderCreateDB, userLogin string,
) {
	res, err := os.request.GetAccrual(order.OrderID)
	if err != nil {
		logger.Log.Error("[StartMonitoringStatus]: GetAccrual error",
			zap.Error(err))
	}

	if res.Status != order.Status {
		order.Status = res.Status
		order.Accrual = res.Accrual

		// TODO: Add transaction for update status + balance
		err := os.db.UpdateOrder(context.Background(),
			order.OrderID, order.Status, order.Accrual)
		if err != nil {
			logger.Log.Error("[StartMonitoringStatus]: update status error", zap.Error(err))
		} else {
			logger.Log.Info("[StartMonitoringStatus]: status updated",
				zap.String("order number", order.OrderID),
				zap.String("prev status", order.Status),
				zap.String("new status", res.Status),
			)
		}
	}

	if order.Status == repository.OrderStatusProcessed {
		os.userDB.UpdateBalance(context.Background(), userLogin, order.Accrual)
		return
	} else if order.Status == repository.OrderStatusInvalid {
		return
	}

	time.Sleep(3 * time.Second)
	os.StartMonitoringStatus(order, userLogin)
}

func (os *OrdersService) UploadOrder(
	ctx context.Context, orderID string,
) error {
	err := os.validator.ValidateVar(orderID, "required,luhn-order")

	if err != nil {
		return err
	}

	user, err := jwt.GetUserFromCtx(ctx)
	if err != nil {
		logger.Log.Error("error while getting user data from context")
		return err
	}

	isErrGetAccrual := false
	res, err := os.request.GetAccrual(orderID)
	if err != nil {
		isErrGetAccrual = true
	}

	order := models.OrderCreateDB{
		OrderID:   orderID,
		UserLogin: user.Login,
	}

	if isErrGetAccrual {
		order.Status = repository.OrderStatusNew
	} else {
		order.Status = res.Status
		order.Accrual = res.Accrual
	}

	err = os.db.Create(ctx, order, user)
	if err != nil {
		return err
	}

	go os.StartMonitoringStatus(order, user.Login)
	return nil
}
