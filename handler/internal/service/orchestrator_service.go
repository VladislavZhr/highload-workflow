package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
)

type OrchestratorService struct {
	db                       *sql.DB
	processingService        *ProcessingService
	finalizeTx2Service       *FinalizeTx2Service
	finalizeFailedTx2Service *FinalizeFailedTx2Service
}

func NewOrchestratorService(
	db *sql.DB,
	processingService *ProcessingService,
	finalizeTx2Service *FinalizeTx2Service,
	finalizeFailedTx2Service *FinalizeFailedTx2Service,
) *OrchestratorService {
	return &OrchestratorService{
		db:                       db,
		processingService:        processingService,
		finalizeTx2Service:       finalizeTx2Service,
		finalizeFailedTx2Service: finalizeFailedTx2Service,
	}
}

func (s *OrchestratorService) HandleMessage(
	ctx context.Context,
	msg []byte,
	now time.Time,
) (ProcessingDecision, error) {
	var transport model.TransportMessage
	if err := decodeTransportMessage(msg, &transport); err != nil {
		return "", err
	}

	tx1, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	decision, err := s.processingService.AcquireForProcessing(
		ctx,
		tx1,
		transport.Message.Header.CorrelationID,
		transport.Message.Header.RequestID,
		now,
	)
	if err != nil {
		_ = tx1.Rollback()
		return "", err
	}

	if err := tx1.Commit(); err != nil {
		return "", err
	}

	switch decision {
	case DecisionSkipCompleted, DecisionSkipPermanent, DecisionSkipProcessing:
		return decision, nil
	case DecisionStartNew, DecisionRetry:
	default:
		return decision, nil
	}

	processed, err := s.processingService.ProcessAfterAcquire(ctx, msg)
	if err != nil {
		tx2, txErr := s.db.BeginTx(ctx, nil)
		if txErr != nil {
			return decision, errors.Join(err, txErr)
		}

		lastError := err.Error()

		if finalizeErr := s.finalizeFailedTx2Service.Finalize(
			ctx,
			tx2,
			transport.Message.Header.CorrelationID,
			model.StatusFailedRetryable,
			&lastError,
		); finalizeErr != nil {
			_ = tx2.Rollback()
			return decision, errors.Join(err, finalizeErr)
		}

		if commitErr := tx2.Commit(); commitErr != nil {
			return decision, errors.Join(err, commitErr)
		}

		return decision, err
	}

	tx2, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return decision, err
	}

	if err := s.finalizeTx2Service.Complete(ctx, tx2, processed); err != nil {
		_ = tx2.Rollback()
		return decision, err
	}

	if err := tx2.Commit(); err != nil {
		return decision, err
	}

	return decision, nil
}

func decodeTransportMessage(msg []byte, transport *model.TransportMessage) error {
	return json.Unmarshal(msg, transport)
}
