package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
	"github.com/VladislavZhr/highload-workflow/handler/internal/repository"
)

type ProcessingService struct {
	processingRepo  *repository.ProcessingRepository
	consumerService *ConsumerService
	leaseDuration   time.Duration
	maxRetryCount   int
}

func NewProcessingService(
	processingRepo *repository.ProcessingRepository,
	consumerService *ConsumerService,
	leaseDuration time.Duration,
	maxRetryCount int,
) *ProcessingService {
	return &ProcessingService{
		processingRepo:  processingRepo,
		consumerService: consumerService,
		leaseDuration:   leaseDuration,
		maxRetryCount:   maxRetryCount,
	}
}

func (s *ProcessingService) AcquireForProcessing(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
	requestID string,
	now time.Time,
) (ProcessingDecision, error) {
	record, err := s.processingRepo.GetByCorrelationID(ctx, tx, correlationID)
	if err != nil {
		return "", err
	}

	decision := DecideProcessingAction(record, now)

	switch decision {
	case DecisionStartNew:
		leaseUntil := now.Add(s.leaseDuration)

		inserted, err := s.processingRepo.InsertProcessing(
			ctx,
			tx,
			model.ProcessingRecord{
				CorrelationID: correlationID,
				RequestID:     requestID,
				Status:        model.StatusProcessing,
				LeaseUntil:    &leaseUntil,
			},
		)
		if err != nil {
			return "", err
		}

		if inserted {
			return DecisionStartNew, nil
		}

		record, err = s.processingRepo.GetByCorrelationID(ctx, tx, correlationID)
		if err != nil {
			return "", err
		}

		return s.resolveRetryDecision(ctx, tx, record, correlationID, requestID, now)

	case DecisionRetry:
		return s.resolveRetryDecision(ctx, tx, record, correlationID, requestID, now)

	case DecisionSkipCompleted, DecisionSkipPermanent, DecisionSkipProcessing:
		return decision, nil

	default:
		return decision, nil
	}
}

func (s *ProcessingService) ProcessAfterAcquire(
	ctx context.Context,
	msg []byte,
) (model.ProcessedMessage, error) {
	return s.consumerService.ProcessMessage(ctx, msg)
}

func (s *ProcessingService) resolveRetryDecision(
	ctx context.Context,
	tx *sql.Tx,
	record *model.ProcessingRecord,
	correlationID string,
	requestID string,
	now time.Time,
) (ProcessingDecision, error) {
	if record == nil {
		return DecisionStartNew, nil
	}

	if record.RetryCount >= s.maxRetryCount {
		marked, err := s.processingRepo.MarkFailedPermanentIfRetryExceeded(
			ctx,
			tx,
			correlationID,
			s.maxRetryCount,
			nil,
		)
		if err != nil {
			return "", err
		}

		if marked {
			return DecisionSkipPermanent, nil
		}

		record, err = s.processingRepo.GetByCorrelationID(ctx, tx, correlationID)
		if err != nil {
			return "", err
		}

		if record == nil {
			return DecisionStartNew, nil
		}

		return DecideProcessingAction(record, now), nil
	}

	leaseUntil := now.Add(s.leaseDuration)

	acquiredExpired, err := s.processingRepo.AcquireExpiredProcessing(
		ctx,
		tx,
		correlationID,
		requestID,
		leaseUntil,
	)
	if err != nil {
		return "", err
	}

	if acquiredExpired {
		return DecisionRetry, nil
	}

	acquiredRetryable, err := s.processingRepo.AcquireRetryableFailed(
		ctx,
		tx,
		correlationID,
		requestID,
		leaseUntil,
	)
	if err != nil {
		return "", err
	}

	if acquiredRetryable {
		return DecisionRetry, nil
	}

	record, err = s.processingRepo.GetByCorrelationID(ctx, tx, correlationID)
	if err != nil {
		return "", err
	}

	if record == nil {
		return DecisionStartNew, nil
	}

	if record.RetryCount >= s.maxRetryCount {
		marked, err := s.processingRepo.MarkFailedPermanentIfRetryExceeded(
			ctx,
			tx,
			correlationID,
			s.maxRetryCount,
			nil,
		)
		if err != nil {
			return "", err
		}

		if marked {
			return DecisionSkipPermanent, nil
		}

		record, err = s.processingRepo.GetByCorrelationID(ctx, tx, correlationID)
		if err != nil {
			return "", err
		}

		if record == nil {
			return DecisionStartNew, nil
		}
	}

	return DecideProcessingAction(record, now), nil
}
