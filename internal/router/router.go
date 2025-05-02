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
		r.Use(middleware.GzipBody)
		r.Use(chiMiddleware.AllowContentEncoding("gzip"))
		r.Use(chiMiddleware.Compress(gzip.BestSpeed))

		r.Route("/api/user", func(r chi.Router) {
			// JSON
			r.Group(func(r chi.Router) {
				r.Use(middleware.AllowJSONOnly)

				r.Post("/register", handlers.User.Register)
				r.Post("/login", handlers.User.Login)
			})

			// with auth
			r.Group(func(r chi.Router) {
				r.Use(middleware.Auth(secret))

				r.Get("/info", handlers.User.Info)
				r.Post("/orders", handlers.Orders.UploadOrder)
				r.Get("/orders", handlers.Orders.GetOrders)
				r.Get("/balance", handlers.User.Balance)
				r.Post("/balance/withdraw", handlers.Withdraw.Withdraw)
				r.Get("/withdrawals", handlers.Withdraw.GetWithdrawals)
			})

			r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				httputils.SuccessJSON(w)
			})
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}
