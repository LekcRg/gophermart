package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type DBUser struct {
	ID           int    `json:"id,omitempty"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
}

type JWTClaim struct {
	jwt.RegisteredClaims
	Login string
	ID    int
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,password,min=8,max=40"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type UserBalance struct {
	Balance   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
