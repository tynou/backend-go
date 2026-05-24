package service

import (
	"context"
	"fmt"
	"payment/internal/producer"
	"payment/internal/repository"
	"pkg/events"
)

type PaymentService struct {
	repo     *repository.PaymentRepository
	producer *producer.PaymentInitProducer
}

func NewPaymentService(repo *repository.PaymentRepository, producer *producer.PaymentInitProducer) *PaymentService {
	return &PaymentService{repo: repo, producer: producer}
}

func (s *PaymentService) Pay(ctx context.Context, userID int32, amount float64) error {
	payment, err := s.repo.CreatePayment(ctx, userID, amount)
	if err != nil {
		return fmt.Errorf("ошибка создания платежа: %w", err)
	}

	return s.producer.Publish(ctx, events.PaymentInit{
		PaymentID: payment.ID,
		UserID:    userID,
		Amount:    amount,
	})
}
