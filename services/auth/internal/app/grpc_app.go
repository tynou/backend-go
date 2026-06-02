package app

import (
	"auth/internal/server"
	"auth/internal/service"
	"context"
	"fmt"
	"log/slog"
	"net"
	"pkg/api/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewGRPCApp(log *slog.Logger, authService *service.AuthService, port int) *GRPCApp {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	authGrpcServer := server.NewAuthGRPCServer(authService)
	auth.RegisterAuthServiceServer(gRPCServer, authGrpcServer)

	return &GRPCApp{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *GRPCApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *GRPCApp) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	a.log.Info("gRPC server is running", slog.Int("port", a.port))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

func (a *GRPCApp) Stop() {
	a.log.Info("stopping gRPC server", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
