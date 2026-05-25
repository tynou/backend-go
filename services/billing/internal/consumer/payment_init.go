package consumer

import (
	"billing/internal/producer"
	"billing/internal/repository"
	"context"
	"encoding/json"
	"log"
	"pkg/events"

	"github.com/segmentio/kafka-go"
)

type PaymentInitConsumer struct {
	repo     *repository.WalletRepository
	producer *producer.PaymentResultProducer
	reader   *kafka.Reader
}

func NewPaymentInitConsumer(brokers []string, repo *repository.WalletRepository, producer *producer.PaymentResultProducer) *PaymentInitConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "payment.init",
		GroupID:     "billing-service-group",
		StartOffset: kafka.FirstOffset,
	})

	return &PaymentInitConsumer{
		repo:     repo,
		producer: producer,
		reader:   reader,
	}
}

func (c *PaymentInitConsumer) Start(ctx context.Context) {
	log.Println("Консюмер payment.init запущен...")

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("ошибка чтения сообщения: %v", err)
			return
		}

		var event events.PaymentInit
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("ошибка демаршалинга: %v", err)
			c.reader.CommitMessages(ctx, msg)
			continue
		}

		err = c.repo.Deduct(ctx, event.UserID, event.Amount)
		c.producer.Publish(ctx, events.PaymentResult{
			PaymentID: event.PaymentID,
			UserID:    event.UserID,
			Amount:    event.Amount,
			Success:   err == nil,
		})
		//if err != nil {
		//	log.Printf("ошибка при вычете средств: %v", err)
		//}

		log.Printf("Платёж успешно обработан.")
		c.reader.CommitMessages(ctx, msg)
	}
}

func (c *PaymentInitConsumer) Close() error {
	return c.reader.Close()
}
