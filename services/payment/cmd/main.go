package main

import (
	"context"
	"errors"
	"log"
	"net"
	"payment/internal/handlers"
	"payment/internal/repository"
	"payment/internal/server"
	"payment/internal/service"
	"pkg/api/payment"
	"pkg/consumer"
	"pkg/producer"

	_ "payment/docs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m, _ := migrate.New("file://db/migrations", "postgres://postgres:1234@localhost:5435/payment?sslmode=disable")
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	pool, err := pgxpool.New(ctx, "postgres://postgres:1234@localhost:5435/payment")
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

	go paymentResultConsumer.Start(ctx)

	svc := service.NewPaymentService(repo, kafkaProducer)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	paymentGrpcServer := server.NewPaymentGRPCServer(svc)
	payment.RegisterPaymentServiceServer(grpcServer, paymentGrpcServer)

	log.Println("gRPC Payment Service запущен на порту 8082...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
