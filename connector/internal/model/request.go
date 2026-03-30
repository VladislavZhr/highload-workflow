package model

import "encoding/json"

type Request struct {
	Meta           Meta           `json:"meta"`
	Counterparties Counterparties `json:"counterparties"`
}

type Meta struct {
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
}

type Counterparties struct {
	Counterparty []Counterparty `json:"counterparty"`
}

type Counterparty struct {
	CounterpartyID string          `json:"CounterpartyID"`
	Data           json.RawMessage `json:"data"`
}
