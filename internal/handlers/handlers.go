package handlers

import (
	"github.com/LekcRg/gophermart/internal/handlers/user"
	"github.com/LekcRg/gophermart/internal/service"
)

type Handlers struct {
	User *user.UserHandler
}

func New(s *service.Service) *Handlers {
	return &Handlers{
		User: user.New(&s.User),
	}
}
