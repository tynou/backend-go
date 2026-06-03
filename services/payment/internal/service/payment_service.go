package service

import (
	"context"
	"log/slog"
	"payment/internal/repository"
	"pkg/events"
	"pkg/producer"
)

type PaymentService struct {
	repo     *repository.PaymentRepository
	producer *producer.KafkaProducer
	log      *slog.Logger
}

func NewPaymentService(repo *repository.PaymentRepository, producer *producer.KafkaProducer, log *slog.Logger) *PaymentService {
	return &PaymentService{repo: repo, producer: producer, log: log}
}

func (s *PaymentService) Pay(ctx context.Context, userID int32, amount float64) error {
	payment, err := s.repo.CreatePayment(ctx, userID, amount)
	if err != nil {
		s.log.Error("failed to create payment", slog.Any("err", err))
		return err
	}

	err = s.producer.Publish(ctx, events.PaymentInit{
		PaymentID: payment.ID,
		UserID:    userID,
		Amount:    amount,
	})
	if err != nil {
		s.log.Error("failed to send payment init event", slog.Any("err", err))
		return err
	}

	return nil
}
