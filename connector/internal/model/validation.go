package model

import "errors"

var (
	ErrMissingCounterparty = errors.New("counterparty array is required")
	ErrEmptyCounterparty   = errors.New("counterparty array is empty")
	ErrNoValidCounterparty = errors.New("no valid counterparty found")
)

func (req *Request) Validate() error {
	if req.Counterparties.Counterparty == nil {
		return ErrMissingCounterparty
	}

	if len(req.Counterparties.Counterparty) == 0 {
		return ErrEmptyCounterparty
	}

	for _, cp := range req.Counterparties.Counterparty {
		if isValidCounterparty(cp) {
			return nil
		}
	}
	return ErrNoValidCounterparty
}

func isValidCounterparty(cp Counterparty) bool {
	if cp.CounterpartyID != "" {
		return true
	}
	if len(cp.Data) > 0 {
		return true
	}
	return false
}
