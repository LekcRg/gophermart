package user

import (
	"context"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/crypto"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db        repository.UserRepository
	validator *validator.Validator
	config    config.Config
}

func New(
	db repository.UserRepository, valid *validator.Validator,
	cfg config.Config,
) *UserService {
	return &UserService{
		db:        db,
		validator: valid,
		config:    cfg,
	}
}

func (us *UserService) Login(
	ctx context.Context, user models.LoginRequest,
) (string, error) {
	err := us.validator.ValidateStruct(user)
	if err != nil {
		return "", err
	}
	dbUser, err := us.db.Login(ctx, user)
	if err != nil {
		return "", err
	}

	return crypto.BuildJWTString(dbUser.ID, us.config.JWTExp, dbUser.Login, us.config.JWTSecret)
}

func (us *UserService) Register(
	ctx context.Context, user models.RegisterRequest,
) (string, error) {
	err := us.validator.ValidateStruct(user)
	if err != nil {
		return "", err
	}
	bPassHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	reqUser := models.DBUser{
		PasswordHash: string(bPassHash),
		Login:        user.Login,
	}
	dbUser, err := us.db.Create(ctx, reqUser)
	if err != nil {
		return "", err
	}

	return crypto.BuildJWTString(dbUser.ID, us.config.JWTExp, dbUser.Login, us.config.JWTSecret)
}

func (us *UserService) Balance(ctx context.Context) (models.UserBalance, error) {
	user, err := crypto.GetUserFromCtx(ctx)
	if err != nil {
		return models.UserBalance{}, err
	}

	return us.db.GetBalance(ctx, user.Login)
}
