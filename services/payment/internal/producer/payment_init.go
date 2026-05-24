package producer

import (
	"context"
	"encoding/json"
	"pkg/events"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type PaymentInitProducer struct {
	writer *kafka.Writer
}

func NewPaymentInitProducer(brokers []string) *PaymentInitProducer {
	return &PaymentInitProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "payment.init",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *PaymentInitProducer) Publish(ctx context.Context, event events.PaymentInit) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(int64(event.UserID), 10)),
		Value: eventBytes,
	})
}

func (p *PaymentInitProducer) Close() error {
	return p.writer.Close()
}
