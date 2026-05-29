package repository

import (
	"auth/internal/db"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{queries: db.New(pool)}
}

func (r *UserRepository) CreateUser(ctx context.Context, username, hash string) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		PasswordHash: hash,
	})
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (db.User, error) {
	return r.queries.GetUserByUsername(ctx, username)
}
