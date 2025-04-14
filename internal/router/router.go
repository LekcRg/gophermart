package router

import (
	"compress/gzip"
	"net/http"

	"github.com/LekcRg/gophermart/internal/handlers"
	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(handlers *handlers.Handlers, secret string) chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
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

				r.Get("/is-auth", handlers.User.IsAuth)
			})
		})
	})

	// Maybe make a separate service for swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return r
}
