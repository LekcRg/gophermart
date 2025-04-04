package router

import (
	"compress/gzip"
	"net/http"

	"github.com/LekcRg/gophermart/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.RequestLogger)
	r.Use(middleware.Compress(gzip.BestSpeed))

	r.Route("/user", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(200)
			w.Write([]byte("Hello world!"))
		})
	})

	return r
}
