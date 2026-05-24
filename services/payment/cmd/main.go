package main

import (
	"context"
	"errors"
	"log"
	"payment/internal/handlers"

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

	h := handlers.NewPaymentHandler()

	r := gin.Default()

	r.POST("/pay", h.Pay)
	r.Run(":8082")
}
