package main

import (
	"auth/internal/handlers"
	"auth/internal/repository"
	"auth/internal/service"
	"context"
	"errors"
	"log"
	"pkg/producer"

	_ "auth/docs"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Auth Service API
// @version         1.0
// @description     Микросервис авторизации.
// @host            localhost:8081
// @BasePath        /
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
	h := handlers.NewAuthHandler(svc)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.Run(":8081")
}
