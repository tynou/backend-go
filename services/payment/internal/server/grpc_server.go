package server

import (
	"context"
	"payment/internal/service"
	"pkg/api/payment"
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
		return nil, err
	}

	return &payment.PaymentResponse{
		Message: "payment pending",
	}, nil
}
