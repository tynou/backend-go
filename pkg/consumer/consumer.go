package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"pkg/events"

	"github.com/segmentio/kafka-go"
)

type HandlerFunc[T events.Event] func(ctx context.Context, event T) error

type KafkaConsumer[T events.Event] struct {
	reader  *kafka.Reader
	log     *slog.Logger
	handler HandlerFunc[T]
}

func NewKafkaConsumer[T events.Event](brokers []string, groupID string, log *slog.Logger, handler HandlerFunc[T]) *KafkaConsumer[T] {
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
		log:     log,
		handler: handler,
	}
}

func (c *KafkaConsumer[T]) Start(ctx context.Context) {
	log := c.log.With(slog.String("topic", c.reader.Config().Topic))
	log.Info("consumer is running")

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Error("failed to read message", slog.Any("err", err))
			return
		}

		var event T
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Error("failed to unmarshal message", slog.Any("err", err))
			c.reader.CommitMessages(ctx, msg)
			continue
		}

		err = c.handler(ctx, event)
		if err != nil {
			log.Error("failed to process message", slog.Any("err", err))
			// continue // TODO: возможно, стоит не коммитить, чтобы попытаться обработать позже
		}

		c.reader.CommitMessages(ctx, msg)
	}
}

func (c *KafkaConsumer[T]) Close() error {
	return c.reader.Close()
}
