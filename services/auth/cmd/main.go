package main

import (
	"context"
	"log"
	"services/auth/internal/handlers"
	"services/auth/internal/repository"
	"services/auth/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	m, _ := migrate.New("file://db/migrations", "postgres://postgres:1234@localhost:5433/auth?sslmode=disable")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	ctx := context.Background()
	pool, _ := pgxpool.New(ctx, "postgres://postgres:1234@localhost:5433/auth")
	defer pool.Close()

	repo := repository.NewUserRepository(pool)
	svc := service.NewAuthService(repo)
	h := handlers.NewAuthHandler(svc)

	r := gin.Default()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.Run(":8081")
}
