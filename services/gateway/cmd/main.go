package main

import (
	"gateway/internal/handler"
	"gateway/internal/middleware"
	"log"
	"pkg/api/auth"
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

	authClient := auth.NewAuthServiceClient(authConn)
	paymentClient := payment.NewPaymentServiceClient(paymentConn)

	authHandler := handler.NewAuthHandler(authClient)
	paymentHandler := handler.NewPaymentHandler(paymentClient)

	r := gin.Default()

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/api/payment/pay", paymentHandler.Pay)
	}

	log.Println("API Gateway запущен на порту 8084...")
	r.Run(":8084")
}
