package handlers

import (
	"billing/internal/repository"
	"context"
	"pkg/events"
	"pkg/producer"
)

type BillingEventHandler struct {
	repo     *repository.WalletRepository
	producer *producer.KafkaProducer
}

func NewBillingEventHandler(repo *repository.WalletRepository, producer *producer.KafkaProducer) *BillingEventHandler {
	return &BillingEventHandler{repo: repo, producer: producer}
}

func (h *BillingEventHandler) OnUserRegistered(ctx context.Context, event events.UserRegistered) error {
	return h.repo.CreateWallet(ctx, event.UserID)
}

func (h *BillingEventHandler) OnPaymentInit(ctx context.Context, event events.PaymentInit) error {
	err := h.repo.Deduct(ctx, event.UserID, event.Amount)
	h.producer.Publish(ctx, events.PaymentResult{
		PaymentID: event.PaymentID,
		UserID:    event.UserID,
		Amount:    event.Amount,
		Success:   err == nil,
	})
	return err
}
