package model

type ProcessedMessage struct {
	RequestID     string
	CorrelationID string
	XMLBody       []byte
}
