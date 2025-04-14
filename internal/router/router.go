package router

import (
	"compress/gzip"
	"net/http"

	"github.com/LekcRg/gophermart/internal/handlers"
	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func New(handlers *handlers.Handlers, secret string) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(middleware.AllowJSONOnly)
	r.Use(middleware.GzipBody)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(gzip.BestSpeed))

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handlers.User.Register)
		r.Post("/login", handlers.User.Login)
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			httputils.SuccessJSON(w)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(secret))

			r.Get("/is-auth", func(w http.ResponseWriter, r *http.Request) {
				httputils.SuccessJSON(w)
			})
		})
	})

	return r
}
