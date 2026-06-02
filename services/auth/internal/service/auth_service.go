package service

import (
	"auth/internal/config"
	"auth/internal/repository"
	"context"
	"errors"
	"log/slog"
	"pkg/producer"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"pkg/events"
)

type AuthService struct {
	repo     *repository.UserRepository
	producer *producer.KafkaProducer
	log      *slog.Logger
	cfg      *config.Config
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func NewAuthService(repo *repository.UserRepository, producer *producer.KafkaProducer, log *slog.Logger, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, producer: producer, log: log, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return err
	}

	user, err := s.repo.CreateUser(ctx, username, string(hash))
	if err != nil {
		s.log.Error("failed to save user", slog.String("error", err.Error()))
		return err
	}

	err = s.producer.Publish(ctx, events.UserRegistered{
		UserID:   user.ID,
		Username: user.Username,
	})
	if err != nil {
		s.log.Error("failed to send user registration event", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		s.log.Error("failed to find user", slog.String("err", err.Error()))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		s.log.Info("invalid credentials", slog.String("err", err.Error()))
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	signedToken, err := token.SignedString(s.cfg.JWTSecret)
	if err != nil {
		s.log.Error("failed to generate token", slog.String("err", err.Error()))
		return "", err
	}

	return signedToken, nil
}
