package server

import (
	"billing/internal/service"
	"context"
	"pkg/api/billing"
)

type BillingGRPCServer struct {
	billing.UnimplementedBillingServiceServer
	svc *service.BillingService
}

func NewBillingGRPCServer(svc *service.BillingService) *BillingGRPCServer {
	return &BillingGRPCServer{svc: svc}
}

func (s *BillingGRPCServer) Deposit(ctx context.Context, req *billing.DepositRequest) (*billing.DepositResponse, error) {
	err := s.svc.Deposit(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, err
	}

	return &billing.DepositResponse{
		Message: "deposited successfully",
	}, nil
}
