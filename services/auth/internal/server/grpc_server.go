package server

import (
	"auth/internal/service"
	"context"
	"pkg/api/auth"
)

type AuthGRPCServer struct {
	auth.UnimplementedAuthServiceServer
	svc *service.AuthService
}

func NewAuthGRPCServer(svc *service.AuthService) *AuthGRPCServer {
	return &AuthGRPCServer{svc: svc}
}

func (s *AuthGRPCServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	err := s.svc.Register(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Message: "user registered successfully",
	}, nil
}

func (s *AuthGRPCServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	token, err := s.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: token,
	}, nil
}
