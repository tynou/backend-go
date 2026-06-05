package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/service"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
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

	brokers := []string{cfg.KafkaBroker}
	kafkaProducer := producer.NewKafkaProducer(brokers)
	defer kafkaProducer.Close()

	repo := repository.NewUserRepository(pool)
	svc := service.NewAuthService(repo, kafkaProducer, log, cfg)

	grpcApp := app.NewGRPCApp(log, svc, cfg.Port)

	go grpcApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	grpcApp.Stop()
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
