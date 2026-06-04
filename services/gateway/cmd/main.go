package main

import (
	"fmt"
	"gateway/internal/config"
	"gateway/internal/handler"
	"gateway/internal/middleware"
	"log/slog"
	"os"
	"pkg/api/auth"
	"pkg/api/billing"
	"pkg/api/payment"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "gateway/docs"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg := config.MustLoad()
	log := setupLogger()

	authConn, err := grpc.NewClient(cfg.AuthAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to auth service", slog.Any("err", err))
		os.Exit(1)
	}
	defer authConn.Close()

	paymentConn, err := grpc.NewClient(cfg.PaymentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to payment service", slog.Any("err", err))
		os.Exit(1)
	}
	defer paymentConn.Close()

	billingConn, err := grpc.NewClient(cfg.BillingAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to billing service", slog.Any("err", err))
		os.Exit(1)
	}
	defer billingConn.Close()

	authClient := auth.NewAuthServiceClient(authConn)
	paymentClient := payment.NewPaymentServiceClient(paymentConn)
	billingClient := billing.NewBillingServiceClient(billingConn)

	authHandler := handler.NewAuthHandler(authClient, log)
	paymentHandler := handler.NewPaymentHandler(paymentClient, log)
	billingHandler := handler.NewBillingHandler(billingClient, log)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("/api/payment/pay", paymentHandler.Pay)

		protected.POST("/api/billing/deposit", billingHandler.Deposit)
	}

	log.Info("gateway is running", slog.Int("port", cfg.Port))

	r.Run(fmt.Sprintf(":%d", cfg.Port))
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
