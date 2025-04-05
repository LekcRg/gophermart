package user

import (
	"context"
	"fmt"

	"github.com/LekcRg/gophermart/internal/models"
	"github.com/LekcRg/gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db repository.UserRepository
}

func New(db repository.UserRepository) *UserService {
	return &UserService{
		db: db,
	}
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (us *UserService) Register(ctx context.Context, user models.AuthRequest) error {
	bPassHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	dbUser := models.User{
		PasswordHash: string(bPassHash),
		Login:        user.Login,
	}
	fmt.Println(string(bPassHash))
	return us.db.Create(ctx, dbUser)
}
