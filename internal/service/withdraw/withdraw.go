package withdraw

import (
	"context"

	"github.com/LekcRg/gophermart/internal/crypto"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/validator"
)

type WithdrawService struct {
	db        repository.WithdrawRepository
	userDB    repository.UserRepository
	validator *validator.Validator
}

func New(
	db repository.WithdrawRepository,
	userDB repository.UserRepository,
	validator *validator.Validator,
) *WithdrawService {

	return &WithdrawService{
		db:        db,
		userDB:    userDB,
		validator: validator,
	}
}

func (ws *WithdrawService) CreateWithdraw(
	ctx context.Context, withdraw models.WithdrawRequest,
) error {
	errs := ws.validator.ValidateStruct(withdraw)
	if errs != nil {
		return errs
	}

	user, err := crypto.GetUserFromCtx(ctx)
	if err != nil {
		return err
	}
	err = ws.userDB.WithdrawBalance(ctx, user.Login, withdraw)
	if err != nil {
		return err
	}

	return ws.db.CreateWithdraw(ctx, user.Login, withdraw)
}

func (ws *WithdrawService) GetWithdrawals(ctx context.Context) (models.WithdrawList, error) {
	user, err := crypto.GetUserFromCtx(ctx)
	if err != nil {
		return models.WithdrawList{}, err
	}

	return ws.db.GetByUserLogin(ctx, user.Login)
}
