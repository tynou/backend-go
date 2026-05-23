package producer

import (
	"context"
	"encoding/json"
	"pkg/events"

	"github.com/segmentio/kafka-go"
)

type UserRegisteredProducer struct {
	writer *kafka.Writer
}

func NewUserRegisteredProducer(brokers []string) *UserRegisteredProducer {
	return &UserRegisteredProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "user.registered",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *UserRegisteredProducer) Publish(ctx context.Context, event events.UserRegistered) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Username),
		Value: eventBytes,
	})
}

func (p *UserRegisteredProducer) Close() error {
	return p.writer.Close()
}
