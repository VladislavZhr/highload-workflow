package model

import "time"

type ProcessingRecord struct {
	CorrelationID string
	RequestID     string
	Status        ProcessingStatus
	LeaseUntil    *time.Time
	RetryCount    int
	LastError     *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
