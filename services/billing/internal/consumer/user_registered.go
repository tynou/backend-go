package consumer

import (
	"billing/internal/repository"
	"context"
	"encoding/json"
	"log"
	"pkg/events"

	"github.com/segmentio/kafka-go"
)

type UserRegisteredConsumer struct {
	repo   *repository.WalletRepository
	reader *kafka.Reader
}

func NewUserRegisteredConsumer(brokers []string, repo *repository.WalletRepository) *UserRegisteredConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "user.registered",
		GroupID:     "billing-service-group",
		StartOffset: kafka.FirstOffset,
	})

	return &UserRegisteredConsumer{
		repo:   repo,
		reader: reader,
	}
}

func (c *UserRegisteredConsumer) Start(ctx context.Context) {
	log.Println("Консюмер user.registered запущен...")

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("ошибка чтения сообщения: %v", err)
			return
		}

		var event events.UserRegistered
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("ошибка демаршалинга: %v", err)
			c.reader.CommitMessages(ctx, msg)
			continue
		}

		err = c.repo.CreateWallet(ctx, event.UserID)
		if err != nil {
			log.Printf("ошибка создания кошелька для %d: %v", event.UserID, err)
			continue
		}

		log.Printf("Кошелек для UserID %d успешно создан!", event.UserID)
		c.reader.CommitMessages(ctx, msg)
	}
}

func (c *UserRegisteredConsumer) Close() error {
	return c.reader.Close()
}
