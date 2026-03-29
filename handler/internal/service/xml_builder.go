package service

import (
	"encoding/xml"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
)

type XMLBuilder struct{}

func NewXMLBuilder() *XMLBuilder {
	return &XMLBuilder{}
}

func (b *XMLBuilder) Build(output model.BusinessOutput) ([]byte, error) {
	data, err := xml.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, err
	}

	result := append([]byte(xml.Header), data...)

	return result, nil
}
