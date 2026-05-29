package repository

import (
	"billing/internal/db"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository struct {
	queries *db.Queries
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{queries: db.New(pool)}
}

func (r *WalletRepository) CreateWallet(ctx context.Context, userID int32) error {
	return r.queries.CreateWallet(ctx, userID)
}

func (r *WalletRepository) Deduct(ctx context.Context, userID int32, amount float64) error {
	var numericAmount pgtype.Numeric
	amountStr := fmt.Sprintf("%.2f", amount)

	err := numericAmount.Scan(amountStr)
	if err != nil {
		return fmt.Errorf("ошибка конвертации суммы: %w", err)
	}

	rowsAffected, err := r.queries.DeductBalance(ctx, db.DeductBalanceParams{
		Balance: numericAmount,
		UserID:  userID,
	})
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("недостаточно средств")
	}

	return nil
}

func (r *WalletRepository) Deposit(ctx context.Context, userID int32, amount float64) error {
	var numericAmount pgtype.Numeric
	amountStr := fmt.Sprintf("%.2f", amount)

	err := numericAmount.Scan(amountStr)
	if err != nil {
		return fmt.Errorf("ошибка конвертации суммы: %w", err)
	}

	return r.queries.DepositBalance(ctx, db.DepositBalanceParams{
		Balance: numericAmount,
		UserID:  userID,
	})
}
