package models

import "time"

type WithdrawRequest struct {
	Order string  `json:"order" validate:"required,luhn-order"`
	Sum   float64 `json:"sum" validate:"required"`
}

type Withdraw struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawList []Withdraw
