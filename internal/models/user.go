package models

type User struct {
	ID           int
	Login        string
	PasswordHash string
}

type AuthRequest struct {
	Login    string
	Password string
}
