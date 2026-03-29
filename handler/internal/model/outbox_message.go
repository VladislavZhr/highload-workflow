package model

type OutboxMessage struct {
	CorrelationID string
	RequestID     string
}
