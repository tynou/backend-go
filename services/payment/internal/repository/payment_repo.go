package repository

import (
	"context"
	"fmt"
	"payment/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	queries *db.Queries
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{queries: db.New(pool)}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, userID int32, amount float64) (db.Payment, error) {
	var numericAmount pgtype.Numeric
	amountStr := fmt.Sprintf("%.2f", amount)

	err := numericAmount.Scan(amountStr)
	if err != nil {
		return db.Payment{}, fmt.Errorf("ошибка конвертации суммы: %w", err)
	}

	return r.queries.CreatePayment(ctx, db.CreatePaymentParams{
		UserID: userID,
		Amount: numericAmount,
		Status: db.PaymentStatusTypePending,
	})
}

func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, success bool) error {
	var status db.PaymentStatusType
	if success {
		status = db.PaymentStatusTypeSuccess
	} else {
		status = db.PaymentStatusTypeFailure
	}

	return r.queries.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		ID:     paymentID,
		Status: status,
	})
}
