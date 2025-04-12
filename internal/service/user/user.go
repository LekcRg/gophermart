package user

import (
	"context"

	"github.com/LekcRg/gophermart/internal/config"
	"github.com/LekcRg/gophermart/internal/crypto"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/repository"
	"github.com/LekcRg/gophermart/internal/validator"
	"go.uber.org/zap"
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
	err = us.db.Login(ctx, user)
	if err != nil {
		return "", err
	}

	return crypto.BuildJWTString(user.Login, us.config.JWTSecret)
}

func (us *UserService) Register(
	ctx context.Context, user models.RegisterRequest,
) (string, error) {
	logger.Log.Info("user",
		zap.String("user.Login", user.Login),
		zap.String("user.Password", user.Password),
	)
	err := us.validator.ValidateStruct(user)
	if err != nil {
		return "", err
	}
	bPassHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	dbUser := models.User{
		PasswordHash: string(bPassHash),
		Login:        user.Login,
	}
	err = us.db.Create(ctx, dbUser)
	if err != nil {
		return "", err
	}

	return crypto.BuildJWTString(dbUser.Login, us.config.JWTSecret)
}
