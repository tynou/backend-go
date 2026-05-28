package consumer

import (
	"context"
	"encoding/json"
	"log"
	"pkg/events"

	"github.com/segmentio/kafka-go"
)

type HandlerFunc[T events.Event] func(ctx context.Context, event T) error

type KafkaConsumer[T events.Event] struct {
	reader  *kafka.Reader
	handler HandlerFunc[T]
}

func NewKafkaConsumer[T events.Event](brokers []string, groupID string, handler HandlerFunc[T]) *KafkaConsumer[T] {
	var zeroEvent T
	topic := zeroEvent.GetTopic()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})

	return &KafkaConsumer[T]{
		reader:  reader,
		handler: handler,
	}
}

func (c *KafkaConsumer[T]) Start(ctx context.Context) {
	log.Printf("Консьюмер для топика [%s] запущен...", c.reader.Config().Topic)

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("ошибка чтения сообщения из %s: %v", c.reader.Config().Topic, err)
			return
		}

		var event T
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("ошибка демаршалинга: %v", err)
			c.reader.CommitMessages(ctx, msg)
			continue
		}

		err = c.handler(ctx, event)
		if err != nil {
			log.Printf("ошибка обработки события в %s: %v", c.reader.Config().Topic, err)
			// continue // возможно, стоит не коммитить, чтобы попытаться обработать позже
		}

		c.reader.CommitMessages(ctx, msg)
	}
}

func (c *KafkaConsumer[T]) Close() error {
	return c.reader.Close()
}
