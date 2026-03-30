package kafka

import (
	"context"
	"encoding/json"

	"github.com/VladislavZhr/highload-workflow/connector/internal/model"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
	}

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) Produce(ctx context.Context, msg model.TransportMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := kafkago.Message{
		Key:   []byte(msg.Message.Header.CorrelationID),
		Value: payload,
		Headers: []kafkago.Header{
			{
				Key:   "requestId",
				Value: []byte(msg.Message.Header.RequestID),
			},
			{
				Key:   "correlationId",
				Value: []byte(msg.Message.Header.CorrelationID),
			},
			{
				Key:   "timestamp",
				Value: []byte(msg.Message.Header.Timestamp),
			},
		},
	}

	return p.writer.WriteMessages(ctx, kafkaMsg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
