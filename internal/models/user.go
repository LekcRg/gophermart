package models

type User struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	PasswordHash string
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,password,min=8,max=25"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
