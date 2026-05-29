package service

import (
	"billing/internal/repository"
	"context"
)

type BillingService struct {
	repo *repository.WalletRepository
}

func NewBillingService(repo *repository.WalletRepository) *BillingService {
	return &BillingService{repo: repo}
}

func (s *BillingService) Deposit(ctx context.Context, userID int32, amount float64) error {
	return s.repo.Deposit(ctx, userID, amount)
}
