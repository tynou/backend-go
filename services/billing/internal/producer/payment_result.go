package producer

import (
	"context"
	"encoding/json"
	"pkg/events"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type PaymentResultProducer struct {
	writer *kafka.Writer
}

func NewPaymentResultProducer(brokers []string) *PaymentResultProducer {
	return &PaymentResultProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "payment.result",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *PaymentResultProducer) Publish(ctx context.Context, event events.PaymentResult) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(int64(event.UserID), 10)),
		Value: eventBytes,
	})
}

func (p *PaymentResultProducer) Close() error {
	return p.writer.Close()
}
