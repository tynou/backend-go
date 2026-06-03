package server

import (
	"auth/internal/service"
	"context"
	"errors"
	"pkg/api/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	// TODO: отправлять понятную ошибку, если ник занят

	return &auth.RegisterResponse{
		Message: "user registered successfully",
	}, nil
}

func (s *AuthGRPCServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	token, err := s.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid username or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &auth.LoginResponse{
		Token: token,
	}, nil
}
