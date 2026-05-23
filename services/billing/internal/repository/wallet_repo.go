package repository

import (
	"billing/internal/db"
	"context"

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
