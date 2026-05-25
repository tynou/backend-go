package events

import "github.com/google/uuid"

type UserRegistered struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
}

type PaymentInit struct {
	PaymentID uuid.UUID `json:"payment_id"`
	UserID    int32     `json:"user_id"`
	Amount    float64   `json:"amount"`
}

type PaymentResult struct {
	PaymentID uuid.UUID `json:"payment_id"`
	UserID    int32     `json:"user_id"`
	Amount    float64   `json:"amount"`
	Success   bool      `json:"success"`
}
