package events

import (
	"strconv"

	"github.com/google/uuid"
)

type UserRegistered struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
}

func (e UserRegistered) GetKey() []byte {
	return []byte(e.Username)
}

func (e UserRegistered) GetTopic() string {
	return "user.registered"
}

type PaymentInit struct {
	PaymentID uuid.UUID `json:"payment_id"`
	UserID    int32     `json:"user_id"`
	Amount    float64   `json:"amount"`
}

func (e PaymentInit) GetKey() []byte {
	return []byte(strconv.FormatInt(int64(e.UserID), 10))
}

func (e PaymentInit) GetTopic() string {
	return "payment.init"
}

type PaymentResult struct {
	PaymentID uuid.UUID `json:"payment_id"`
	UserID    int32     `json:"user_id"`
	Amount    float64   `json:"amount"`
	Success   bool      `json:"success"`
}

func (e PaymentResult) GetKey() []byte {
	return []byte(strconv.FormatInt(int64(e.UserID), 10))
}

func (e PaymentResult) GetTopic() string {
	return "payment.result"
}
