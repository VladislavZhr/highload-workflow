package mapper

type ProcessedMessage struct {
	RequestID     string
	CorrelationID string
	XMLBody       []byte
}
