package main

import (
	"context"
	"errors"
	"log"
	"payment/internal/handlers"
	"payment/internal/producer"
	"payment/internal/repository"
	"payment/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

	brokers := []string{"localhost:9092"}
	paymentProducer := producer.NewPaymentInitProducer(brokers)
	defer paymentProducer.Close()

	repo := repository.NewPaymentRepository(pool)
	svc := service.NewPaymentService(repo, paymentProducer)
	h := handlers.NewPaymentHandler(svc)

	r := gin.Default()

	r.POST("/pay", h.Pay)
	r.Run(":8082")
}
