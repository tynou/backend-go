package main

import (
	"gateway/internal/handler"
	"log"
	"pkg/api/auth"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth service: %v", err)
	}
	defer conn.Close()

	authClient := auth.NewAuthServiceClient(conn)
	authHandler := handler.NewAuthHandler(authClient)

	r := gin.Default()

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	log.Println("API Gateway запущен на порту 8084...")
	r.Run(":8084")
}
