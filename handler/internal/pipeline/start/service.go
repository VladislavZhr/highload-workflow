package start

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrSkipCompleted  = errors.New("message already completed")
	ErrSkipPermanent  = errors.New("message permanently failed")
	ErrSkipProcessing = errors.New("message is already being processed")
	ErrStateConflict  = errors.New("message state changed concurrently")
)

type Service struct {
	repo          *Repository
	maxRetryCount int
	leaseDuration time.Duration
}

func NewService(repo *Repository, maxRetryCount int, leaseDuration time.Duration) *Service {
	return &Service{
		repo:          repo,
		maxRetryCount: maxRetryCount,
		leaseDuration: leaseDuration,
	}
}

func (s *Service) InitPipeline(ctx context.Context, tx *sql.Tx, correlationID string, requestID string) error {
	record, err := s.repo.GetByCorrelationID(ctx, tx, correlationID)

	if err != nil {
		return err
	}

	now := time.Now()
	leaseUntil := now.Add(s.leaseDuration)

	if record == nil {
		inserted, err := s.repo.InsertProcessing(ctx, tx, MessageState{
			CorrelationID: correlationID,
			RequestID:     requestID,
			Status:        StatusProcessing,
			LeaseUntil:    &leaseUntil,
		})
		if err != nil {
			return err
		}
		if inserted {
			return nil
		}

		return ErrStateConflict
	}

	switch record.Status {
	case StatusCompleted:
		return ErrSkipCompleted

	case StatusFailedPermanent:
		return ErrSkipPermanent

	case StatusFailedRetryable:

		markedPermanent, err := s.repo.MarkFailedPermanentIfRetryExceeded(
			ctx,
			tx,
			correlationID,
			s.maxRetryCount,
			record.LastError,
		)
		if err != nil {
			return err
		}
		if markedPermanent {
			return ErrSkipPermanent
		}

		retried, err := s.repo.TryRetryFailed(
			ctx,
			tx,
			correlationID,
			s.maxRetryCount,
		)
		if err != nil {
			return err
		}
		if retried {
			return nil
		}

		return ErrStateConflict

	case StatusProcessing:
		if record.LeaseUntil == nil || record.LeaseUntil.Before(now) {
			reacquired, err := s.repo.TryUpdateExpiredLease(
				ctx,
				tx,
				correlationID,
				leaseUntil,
			)
			if err != nil {
				return err
			}
			if reacquired {
				return nil
			}

			return ErrStateConflict
		}

		return ErrSkipProcessing

	default:
		return ErrStateConflict
	}

}
