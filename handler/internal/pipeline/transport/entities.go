package transport

import "encoding/json"

type TransportMessage struct {
	Message TransportEnvelope `json:"message"`
}

type TransportEnvelope struct {
	Header TransportHeader `json:"header"`
	Body   TransportBody   `json:"body"`
}

type TransportHeader struct {
	RequestID     string `json:"requestId"`
	CorrelationID string `json:"correlationId"`
	Timestamp     string `json:"timestamp"`
}

type TransportBody struct {
	Raw json.RawMessage `json:"raw"`
}
