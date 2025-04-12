package router

import (
	"compress/gzip"

	"github.com/LekcRg/gophermart/internal/handlers"
	"github.com/LekcRg/gophermart/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func New(handlers *handlers.Handlers) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(middleware.AllowJSONOnly)
	r.Use(middleware.GzipBody)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(gzip.BestSpeed))

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handlers.User.Register)
		r.Post("/login", handlers.User.Login)
	})

	return r
}
