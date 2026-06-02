package main

import (
	"billing/internal/app"
	"billing/internal/config"
	"billing/internal/handlers"
	"billing/internal/repository"
	"billing/internal/service"
	"context"
	"errors"
	"log/slog"
	"os"
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
	cfg := config.MustLoad()
	log := setupLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m, _ := migrate.New("file://db/migrations", cfg.DBConn)
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error("migration error", slog.Any("err", err))
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, cfg.DBConn)
	if err != nil {
		log.Error("db connection error", slog.Any("err", err))
		os.Exit(1)
	}
	defer pool.Close()

	repo := repository.NewWalletRepository(pool)

	brokers := []string{cfg.KafkaBroker}
	kafkaProducer := producer.NewKafkaProducer(brokers)
	defer kafkaProducer.Close()

	eventHandler := handlers.NewBillingEventHandler(repo, kafkaProducer)

	userRegisteredConsumer := consumer.NewKafkaConsumer(
		brokers,
		"billing-service-group",
		log,
		eventHandler.OnUserRegistered,
	)
	defer userRegisteredConsumer.Close()

	paymentInitConsumer := consumer.NewKafkaConsumer(
		brokers,
		"billing-service-group",
		log,
		eventHandler.OnPaymentInit,
	)
	defer paymentInitConsumer.Close()

	go userRegisteredConsumer.Start(ctx)
	go paymentInitConsumer.Start(ctx)

	svc := service.NewBillingService(repo, log)

	grpcApp := app.NewGRPCApp(log, svc, cfg.Port)

	go grpcApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	grpcApp.Stop()
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
