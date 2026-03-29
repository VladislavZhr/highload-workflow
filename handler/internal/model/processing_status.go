package model

type ProcessingStatus string

const (
	StatusProcessing      ProcessingStatus = "processing"
	StatusCompleted       ProcessingStatus = "completed"
	StatusFailedPermanent ProcessingStatus = "failed_permanent"
	StatusFailedRetryable ProcessingStatus = "failed_retryable"
)
