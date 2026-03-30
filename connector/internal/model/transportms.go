package model

type TransportMessage struct {
	Message Message `json:"message"`
}

type Message struct {
	Header TransportHeader `json:"header"`
	Body   TransportBody   `json:"body"`
}

type TransportHeader struct {
	RequestID     string `json:"requestId"`
	CorrelationID string `json:"correlationId"`
	Timestamp     string `json:"timestamp"`
}

type TransportBody struct {
	Raw BusinessInput `json:"raw"`
}

type BusinessInput struct {
	Counterparties CounterpartiesPayload `json:"counterparties"`
}

type CounterpartiesPayload struct {
	Counterparty []Counterparty `json:"counterparty"`
}
