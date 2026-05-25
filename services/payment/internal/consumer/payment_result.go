package consumer

import (
	"context"
	"encoding/json"
	"log"
	"payment/internal/repository"
	"pkg/events"

	"github.com/segmentio/kafka-go"
)

type PaymentResultConsumer struct {
	repo   *repository.PaymentRepository
	reader *kafka.Reader
}

func NewPaymentResultConsumer(brokers []string, repo *repository.PaymentRepository) *PaymentResultConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "payment.result",
		GroupID:     "payment-service-group",
		StartOffset: kafka.FirstOffset,
	})

	return &PaymentResultConsumer{
		repo:   repo,
		reader: reader,
	}
}

func (c *PaymentResultConsumer) Start(ctx context.Context) {
	log.Println("Консюмер payment.result запущен...")

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("ошибка чтения сообщения: %v", err)
			return
		}

		var event events.PaymentResult
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("ошибка демаршалинга: %v", err)
			c.reader.CommitMessages(ctx, msg)
			continue
		}

		err = c.repo.UpdatePaymentStatus(ctx, event.PaymentID, event.Success)
		if err != nil {
			log.Printf("ошибка обновления статуса платежа: %v", err)
			continue
		}

		log.Printf("Статус платежа успешно изменён.")
		c.reader.CommitMessages(ctx, msg)
	}
}

func (c *PaymentResultConsumer) Close() error {
	return c.reader.Close()
}
