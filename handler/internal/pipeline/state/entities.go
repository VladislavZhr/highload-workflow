package state

import "time"

type ProcessingStatus string

const (
	StatusProcessing      ProcessingStatus = "processing"
	StatusCompleted       ProcessingStatus = "completed"
	StatusFailedPermanent ProcessingStatus = "failed_permanent"
	StatusFailedRetryable ProcessingStatus = "failed_retryable"
)

type MessageState struct {
	CorrelationID string
	RequestID     string
	Status        ProcessingStatus
	LeaseUntil    *time.Time
	RetryCount    int
	LastError     *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
