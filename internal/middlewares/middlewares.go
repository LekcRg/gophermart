package middlewares

import (
	"net/http"
	"time"

	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		duration := time.Since(now)

		logger.Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", ww.Status()),
			zap.Int("size", ww.BytesWritten()),
			zap.String("Content-Encoding", ww.Header().Get("Content-Encoding")),
			zap.Duration("duration", duration),
		)
	})
}
