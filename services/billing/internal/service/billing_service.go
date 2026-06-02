package service

import (
	"billing/internal/repository"
	"context"
	"log/slog"
)

type BillingService struct {
	repo *repository.WalletRepository
	log  *slog.Logger
}

func NewBillingService(repo *repository.WalletRepository, log *slog.Logger) *BillingService {
	return &BillingService{repo: repo, log: log}
}

func (s *BillingService) Deposit(ctx context.Context, userID int32, amount float64) error {
	err := s.repo.Deposit(ctx, userID, amount)
	if err != nil {
		s.log.Error("failed to deposit", slog.Any("err", err))
		return err
	}
	return nil
}
