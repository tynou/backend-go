package main

import (
	"billing/internal/consumer"
	"billing/internal/repository"
	"context"
	"errors"
	"log"
	"os/signal"
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

	userRegisteredConsumer := consumer.NewUserRegisteredConsumer(brokers, repo)
	defer userRegisteredConsumer.Close()
	paymentInitConsumer := consumer.NewPaymentInitConsumer(brokers, repo, kafkaProducer)
	defer paymentInitConsumer.Close()

	go userRegisteredConsumer.Start(ctx)
	go paymentInitConsumer.Start(ctx)

	<-ctx.Done()
}
