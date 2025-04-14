package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"time"

	"github.com/LekcRg/gophermart/internal/httputils"
	"github.com/LekcRg/gophermart/internal/logger"
	"github.com/LekcRg/gophermart/internal/models"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ww := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		duration := time.Since(now)

		logger.Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", ww.Status()),
			zap.Int("size", ww.BytesWritten()),
			zap.String("Content-Type", r.Header.Get("Content-Type")),
			zap.Duration("duration", duration),
		)
	})
}

func GzipBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Log.Error("Error while create gzip reader",
					zap.Error(err))
			}

			r.Body = gz
			defer gz.Close()
		}

		next.ServeHTTP(w, r)
	})
}

func AllowJSONOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			httputils.ErrJSON(w, "Incorrect Content-type", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			token = strings.Replace(token, "Bearer", "", 1)
			token = strings.TrimSpace(token)

			if token == "" {
				logger.Log.Info("[auth middleware]: token is empty")
				httputils.ErrJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var test models.JWTClaim
			parsedToken, err := jwt.ParseWithClaims(token, &test, func(token *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err != nil || !parsedToken.Valid {
				logger.Log.Error("error parse jwt token",
					zap.Error(err))
				httputils.ErrJSON(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			logger.Log.Info("success auth",
				zap.String("login", test.Login),
				zap.String("id", test.Id),
			)

			next.ServeHTTP(w, r)
		})
	}
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
// eyJleHAiOjE3NDQ2NDIwMzgsIkxvZ2luIjoidXNlcjYifQ.
// OA5PHfPt5j5GLi7UrLz6weU7nVXsEoYxnWlitq8bKx0
