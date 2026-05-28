package main

import (
	"billing/internal/handlers"
	"billing/internal/repository"
	"context"
	"errors"
	"log"
	"os/signal"
	"pkg/consumer"
	"pkg/producer"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	m, _ := migrate.New("file://db/migrations", "postgres://postgres:1234@localhost:5434/billing?sslmode=disable")
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	pool, err := pgxpool.New(ctx, "postgres://postgres:1234@localhost:5434/billing")
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	repo := repository.NewWalletRepository(pool)

	brokers := []string{"localhost:9092"}
	kafkaProducer := producer.NewKafkaProducer(brokers)
	defer kafkaProducer.Close()

	eventHandler := handlers.NewBillingEventHandler(repo, kafkaProducer)

	userRegisteredConsumer := consumer.NewKafkaConsumer(
		brokers,
		"billing-service-group",
		eventHandler.OnUserRegistered,
	)
	defer userRegisteredConsumer.Close()
	paymentInitConsumer := consumer.NewKafkaConsumer(
		brokers,
		"billing-service-group",
		eventHandler.OnPaymentInit,
	)
	defer paymentInitConsumer.Close()

	go userRegisteredConsumer.Start(ctx)
	go paymentInitConsumer.Start(ctx)

	<-ctx.Done()
}
