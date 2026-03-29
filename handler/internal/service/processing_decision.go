package service

import (
	"time"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
)

type ProcessingDecision string

const (
	DecisionStartNew       ProcessingDecision = "start_new"
	DecisionRetry          ProcessingDecision = "retry"
	DecisionSkipCompleted  ProcessingDecision = "skip_completed"
	DecisionSkipPermanent  ProcessingDecision = "skip_permanent"
	DecisionSkipProcessing ProcessingDecision = "skip_processing"
)

func DecideProcessingAction(
	record *model.ProcessingRecord,
	now time.Time,
) ProcessingDecision {
	if record == nil {
		return DecisionStartNew
	}

	switch record.Status {
	case model.StatusCompleted:
		return DecisionSkipCompleted

	case model.StatusFailedPermanent:
		return DecisionSkipPermanent

	case model.StatusFailedRetryable:
		return DecisionRetry

	case model.StatusProcessing:
		if record.LeaseUntil == nil {
			return DecisionRetry
		}

		if record.LeaseUntil.Before(now) {
			return DecisionRetry
		}

		return DecisionSkipProcessing

	default:
		return DecisionRetry
	}
}
