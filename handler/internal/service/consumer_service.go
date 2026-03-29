package service

import (
	"context"
	"encoding/json"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
)

type ConsumerService struct {
	mapper     *Mapper
	xmlBuilder *XMLBuilder
}

func NewConsumerService(
	mapper *Mapper,
	xmlBuilder *XMLBuilder,
) *ConsumerService {
	return &ConsumerService{
		mapper:     mapper,
		xmlBuilder: xmlBuilder,
	}
}

func (s *ConsumerService) ProcessMessage(
	ctx context.Context,
	msg []byte,
) (model.ProcessedMessage, error) {
	_ = ctx

	var transport model.TransportMessage
	if err := json.Unmarshal(msg, &transport); err != nil {
		return model.ProcessedMessage{}, err
	}

	var input model.BusinessInput
	if err := json.Unmarshal(transport.Message.Body.Raw, &input); err != nil {
		return model.ProcessedMessage{}, err
	}

	output := s.mapper.Map(input, transport.Message.Header.Timestamp)

	xmlBytes, err := s.xmlBuilder.Build(output)
	if err != nil {
		return model.ProcessedMessage{}, err
	}

	return model.ProcessedMessage{
		RequestID:     transport.Message.Header.RequestID,
		CorrelationID: transport.Message.Header.CorrelationID,
		XMLBody:       xmlBytes,
	}, nil
}
