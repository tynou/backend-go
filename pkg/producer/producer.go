package producer

import (
	"context"
	"encoding/json"
	"pkg/events"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, event events.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: event.GetTopic(),
		Key:   event.GetKey(),
		Value: eventBytes,
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
