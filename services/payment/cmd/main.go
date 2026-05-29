package main

import (
	"context"
	"errors"
	"log"
	"payment/internal/handlers"
	"payment/internal/repository"
	"payment/internal/service"
	"pkg/consumer"
	"pkg/producer"

	_ "payment/docs"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Payment Service API
// @version         1.0
// @description     Микросервис оплаты.
// @host            localhost:8082
// @BasePath        /
func main() {
	m, _ := migrate.New("file://db/migrations", "postgres://postgres:1234@localhost:5435/payment?sslmode=disable")
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), "postgres://postgres:1234@localhost:5435/payment")
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	repo := repository.NewPaymentRepository(pool)
	eventHandler := handlers.NewPaymentEventHandler(repo)

	brokers := []string{"localhost:9092"}
	kafkaProducer := producer.NewKafkaProducer(brokers)
	defer kafkaProducer.Close()

	paymentResultConsumer := consumer.NewKafkaConsumer(
		brokers,
		"payment-service-group",
		eventHandler.OnPaymentResult,
	)
	defer paymentResultConsumer.Close()

	go paymentResultConsumer.Start(context.Background())

	svc := service.NewPaymentService(repo, kafkaProducer)
	h := handlers.NewPaymentHandler(svc)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/pay", h.Pay)
	r.Run(":8082")
}
