package main

import (
	"gateway/internal/handler"
	"gateway/internal/middleware"
	"log"
	"pkg/api/auth"
	"pkg/api/billing"
	"pkg/api/payment"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	authConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	paymentConn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}
	defer paymentConn.Close()

	billingConn, err := grpc.NewClient("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to billing service: %v", err)
	}
	defer billingConn.Close()

	authClient := auth.NewAuthServiceClient(authConn)
	paymentClient := payment.NewPaymentServiceClient(paymentConn)
	billingClient := billing.NewBillingServiceClient(billingConn)

	authHandler := handler.NewAuthHandler(authClient)
	paymentHandler := handler.NewPaymentHandler(paymentClient)
	billingHandler := handler.NewBillingHandler(billingClient)

	r := gin.Default()

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/api/payment/pay", paymentHandler.Pay)

		protected.POST("/api/billing/deposit", billingHandler.Deposit)
	}

	log.Println("API Gateway запущен на порту 8084...")
	r.Run(":8084")
}
