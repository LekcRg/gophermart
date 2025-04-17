package models

import "time"

type OrderDB struct {
	ID         int
	OrderID    string
	Status     string
	UserLogin  string
	UploadedAt time.Time
	UpdatedAt  time.Time
}
