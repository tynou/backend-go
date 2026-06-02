package server

import (
	"billing/internal/service"
	"context"
	"pkg/api/billing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.Internal, "failed to deposit")
	}

	return &billing.DepositResponse{
		Message: "deposited successfully",
	}, nil
}
