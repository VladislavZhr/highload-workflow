package mapper

import (
	"context"
	"encoding/json"
	"encoding/xml"

	"github.com/VladislavZhr/highload-workflow/handler/internal/transport"
)

type Service struct {
	mapper *Mapper
}

func NewService(mapper *Mapper) *Service {
	return &Service{
		mapper: mapper,
	}
}

func (s *Service) ProcessMessage(
	ctx context.Context,
	tm transport.TransportMessage,
) (ProcessedMessage, error) {
	_ = ctx

	var input BusinessInput
	if err := json.Unmarshal(tm.Message.Body.Raw, &input); err != nil {
		return ProcessedMessage{}, err
	}

	output := s.mapper.Map(input, tm.Message.Header.Timestamp)

	xmlBytes, err := buildXML(output)
	if err != nil {
		return ProcessedMessage{}, err
	}

	return ProcessedMessage{
		RequestID:     tm.Message.Header.RequestID,
		CorrelationID: tm.Message.Header.CorrelationID,
		XMLBody:       xmlBytes,
	}, nil
}

func buildXML(output BusinessOutput) ([]byte, error) {
	data, err := xml.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), data...), nil
}
