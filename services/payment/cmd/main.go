package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"payment/internal/app"
	"payment/internal/config"
	"payment/internal/handlers"
	"payment/internal/repository"
	"payment/internal/service"
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
	//defer cancel()

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

	repo := repository.NewPaymentRepository(pool)
	eventHandler := handlers.NewPaymentEventHandler(repo)

	brokers := []string{cfg.KafkaBroker}
	kafkaProducer := producer.NewKafkaProducer(brokers)
	defer kafkaProducer.Close()

	paymentResultConsumer := consumer.NewKafkaConsumer(
		brokers,
		"payment-result-group",
		log,
		eventHandler.OnPaymentResult,
	)
	defer paymentResultConsumer.Close()

	go paymentResultConsumer.Start(ctx)

	svc := service.NewPaymentService(repo, kafkaProducer, log)

	grpcApp := app.NewGRPCApp(log, svc, cfg.Port)

	go grpcApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	cancel()

	grpcApp.Stop()
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
