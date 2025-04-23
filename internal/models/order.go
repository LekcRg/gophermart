package models

import "time"

type OrderDB struct {
	ID         int       `json:"-"`
	OrderID    string    `json:"number"`
	Accrual    float64   `json:"accrual,omitempty"`
	Status     string    `json:"status"`
	UserLogin  string    `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
	UpdatedAt  time.Time `json:"-"`
}

type OrderCreateDB struct {
	OrderID   string
	Status    string
	UserLogin string
	Accrual   float64
}

type AccrualRes struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
