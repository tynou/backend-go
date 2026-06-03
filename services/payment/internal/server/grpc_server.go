package server

import (
	"context"
	"payment/internal/service"
	"pkg/api/payment"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentGRPCServer struct {
	payment.UnimplementedPaymentServiceServer
	svc *service.PaymentService
}

func NewPaymentGRPCServer(svc *service.PaymentService) *PaymentGRPCServer {
	return &PaymentGRPCServer{svc: svc}
}

func (s *PaymentGRPCServer) Pay(ctx context.Context, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	err := s.svc.Pay(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create payment")
	}

	return &payment.PaymentResponse{
		Message: "payment pending",
	}, nil
}
