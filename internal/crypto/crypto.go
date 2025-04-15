package crypto

import (
	"context"
	"fmt"
	"time"

	"github.com/LekcRg/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type key int

const UserContextKey key = iota

func BuildJWTString(id int, exp time.Duration, login, secret string) (string, error) {
	fmt.Println(jwt.NewNumericDate(time.Now().Add(2 * time.Hour)))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Login: login,
		ID:    id,
	})

	return token.SignedString([]byte(secret))
}

func GetUserClaims(token, secret string) (models.JWTClaim, error) {
	var claim models.JWTClaim
	parsedToken, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return claim, err
	}

	if !parsedToken.Valid {
		return claim, fmt.Errorf("not valid token")
	}

	return claim, nil
}

func GetUserFromCtx(ctx context.Context) (models.DBUser, error) {
	user, ok := ctx.Value(UserContextKey).(models.DBUser)
	if !ok {
		return models.DBUser{}, fmt.Errorf("not found")
	}

	return user, nil
}
