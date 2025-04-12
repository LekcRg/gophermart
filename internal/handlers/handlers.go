package handlers

import (
	"github.com/LekcRg/gophermart/internal/handlers/user"
	"github.com/LekcRg/gophermart/internal/service"
	"github.com/LekcRg/gophermart/internal/validator"
)

type Handlers struct {
	User *user.UserHandler
}

func New(s *service.Service, validator *validator.Validator) *Handlers {
	return &Handlers{
		User: user.New(&s.User, validator),
	}
}
