package crypto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TOKEN_EXP = time.Hour * 3

func BuildJWTString(login string, secret string) (string, error) {
	type Claims struct {
		jwt.RegisteredClaims
		Login string
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		Login: login,
	})

	return token.SignedString([]byte(secret))
}
