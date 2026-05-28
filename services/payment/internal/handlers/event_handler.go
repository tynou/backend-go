package handlers

import (
	"context"
	"payment/internal/repository"
	"pkg/events"
)

type PaymentEventHandler struct {
	repo *repository.PaymentRepository
}

func NewPaymentEventHandler(repo *repository.PaymentRepository) *PaymentEventHandler {
	return &PaymentEventHandler{repo: repo}
}

func (h *PaymentEventHandler) OnPaymentResult(ctx context.Context, event events.PaymentResult) error {
	return h.repo.UpdatePaymentStatus(ctx, event.PaymentID, event.Success)
}
