package crypto

import (
	"time"

	"github.com/LekcRg/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const TOKEN_EXP = time.Hour * 3

func BuildJWTString(login string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		Login: login,
	})

	return token.SignedString([]byte(secret))
}
