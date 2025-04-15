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

// TODO: get token_exp from config
const TOKEN_EXP = time.Hour * 3

func BuildJWTString(id int, login, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
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

func GetUserFromCtx(ctx context.Context) (models.JWTClaim, error) {
	user, ok := ctx.Value(UserContextKey).(models.JWTClaim)
	if !ok {
		return models.JWTClaim{}, fmt.Errorf("not found")
	}

	return user, nil
}
