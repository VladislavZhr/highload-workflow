package service

import (
	"context"
	"errors"

	"github.com/VladislavZhr/highload-workflow/connector/internal/kafka"
	"github.com/VladislavZhr/highload-workflow/connector/internal/model"
	"github.com/google/uuid"
)

var (
	ErrValidation   = errors.New("validation error")
	ErrKafkaProduce = errors.New("kafka produce error")
)

type ConnectorService struct {
	producer *kafka.Producer
}

func NewConnectorService(producer *kafka.Producer) *ConnectorService {
	return &ConnectorService{
		producer: producer,
	}
}

func (s *ConnectorService) Process(ctx context.Context, req model.Request) (model.TransportMessage, error) {

	if err := req.Validate(); err != nil {
		return model.TransportMessage{}, errors.Join(ErrValidation, err)
	}

	transport := model.TransportMessage{
		Message: model.Message{
			Header: model.TransportHeader{
				RequestID:     req.Meta.RequestID,
				CorrelationID: uuid.NewString(),
				Timestamp:     req.Meta.Timestamp,
			},
			Body: model.TransportBody{
				Raw: req.Counterparties.Counterparty,
			},
		},
	}

	if err := s.producer.Produce(ctx, transport); err != nil {
		return model.TransportMessage{}, errors.Join(ErrKafkaProduce, err)
	}

	return transport, nil
}
