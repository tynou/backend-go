package service

import (
	"context"
	"errors"
	"pkg/producer"
	"services/auth/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"pkg/events"
)

var jwtSecret = []byte("my_super_secret_key") // TODO: сделать по-нормальному

type AuthService struct {
	repo     *repository.UserRepository
	producer *producer.KafkaProducer
}

func NewAuthService(repo *repository.UserRepository, producer *producer.KafkaProducer) *AuthService {
	return &AuthService{repo: repo, producer: producer}
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user, err := s.repo.CreateUser(ctx, username, string(hash))
	if err != nil {
		return err
	}

	return s.producer.Publish(ctx, events.UserRegistered{
		UserID:   user.ID,
		Username: user.Username,
	})
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(jwtSecret)
}
