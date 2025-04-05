package router

import (
	"compress/gzip"
	"net/http"

	"github.com/LekcRg/gophermart/internal/handlers"
	"github.com/LekcRg/gophermart/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func New(handlers *handlers.Handlers) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(chiMiddleware.Compress(gzip.BestSpeed))

	r.Route("/api/user", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(200)
			w.Write([]byte("Hello world!"))
		})
		r.Post("/register", handlers.User.Register)
	})

	return r
}
