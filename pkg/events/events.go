package events

type UserRegistered struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
}

type PaymentInit struct {
	PaymentID string  `json:"payment_id"`
	UserID    int32   `json:"user_id"`
	Amount    float64 `json:"amount"`
}

type PaymentResult struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}
