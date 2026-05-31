package main

import (
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"context"
	"errors"
	"log"
	"net"
	"pkg/api/auth"
	"pkg/producer"

	_ "auth/docs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m, _ := migrate.New("file://db/migrations", "postgres://postgres:1234@localhost:5433/auth?sslmode=disable")
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	pool, err := pgxpool.New(ctx, "postgres://postgres:1234@localhost:5433/auth")
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	brokers := []string{"localhost:9092"}
	kafkaProducer := producer.NewKafkaProducer(brokers)
	defer kafkaProducer.Close()

	repo := repository.NewUserRepository(pool)
	svc := service.NewAuthService(repo, kafkaProducer)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	authGrpcServer := server.NewAuthGRPCServer(svc)
	auth.RegisterAuthServiceServer(grpcServer, authGrpcServer)

	log.Println("gRPC Auth Service запущен на порту 8081...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
