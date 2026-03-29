package pipeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/finalize"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/mapper"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/start"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/state"
	"github.com/VladislavZhr/highload-workflow/handler/internal/transport"
)

type Orchestrator struct {
	db              *sql.DB
	startService    *start.Service
	mapperService   *mapper.Service
	finalizeService *finalize.Service
}

func NewOrchestrator(
	db *sql.DB,
	startService *start.Service,
	mapperService *mapper.Service,
	finalizeService *finalize.Service,
) *Orchestrator {
	return &Orchestrator{
		db:              db,
		startService:    startService,
		mapperService:   mapperService,
		finalizeService: finalizeService,
	}
}

func (o *Orchestrator) Handle(
	ctx context.Context,
	tm transport.TransportMessage,
) (mapper.ProcessedMessage, error) {
	correlationID := tm.Message.Header.CorrelationID
	requestID := tm.Message.Header.RequestID

	if err := o.runStartTx(ctx, correlationID, requestID); err != nil {
		return mapper.ProcessedMessage{}, err
	}

	processedMsg, processErr := o.mapperService.ProcessMessage(ctx, tm)
	if processErr != nil {
		status := classifyProcessingError(processErr)
		lastError := stringPtr(processErr.Error())

		finalizeErr := o.runFinalizeFailedTx(
			ctx,
			correlationID,
			status,
			lastError,
		)
		if finalizeErr != nil {
			return mapper.ProcessedMessage{}, fmt.Errorf(
				"processing failed: %v; finalize failed tx also failed: %w",
				processErr,
				finalizeErr,
			)
		}

		return mapper.ProcessedMessage{}, processErr
	}

	if err := o.runFinalizeSuccessTx(ctx, correlationID); err != nil {
		return mapper.ProcessedMessage{}, err
	}

	return processedMsg, nil
}

func (o *Orchestrator) runStartTx(
	ctx context.Context,
	correlationID string,
	requestID string,
) error {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin start tx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := o.startService.InitPipeline(ctx, tx, correlationID, requestID); err != nil {
		if errors.Is(err, start.ErrSkipCompleted) ||
			errors.Is(err, start.ErrSkipPermanent) ||
			errors.Is(err, start.ErrSkipProcessing) ||
			errors.Is(err, start.ErrStateConflict) {
			return err
		}

		return fmt.Errorf("init pipeline: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit start tx: %w", err)
	}

	committed = true
	return nil
}

func (o *Orchestrator) runFinalizeSuccessTx(
	ctx context.Context,
	correlationID string,
) error {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin finalize success tx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := o.finalizeService.FinalizeSuccessState(ctx, tx, correlationID); err != nil {
		return fmt.Errorf("finalize success state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit finalize success tx: %w", err)
	}

	committed = true
	return nil
}

func (o *Orchestrator) runFinalizeFailedTx(
	ctx context.Context,
	correlationID string,
	status state.ProcessingStatus,
	lastError *string,
) error {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin finalize failed tx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := o.finalizeService.FinalizeFailedState(
		ctx,
		tx,
		correlationID,
		status,
		lastError,
	); err != nil {
		return fmt.Errorf("finalize failed state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit finalize failed tx: %w", err)
	}

	committed = true
	return nil
}

func classifyProcessingError(err error) state.ProcessingStatus {
	// Поки що всі помилки обробки вважаємо retryable.
	// Потім сюди можна додати класифікацію:
	// malformed JSON -> failed_permanent
	// transient infra issue -> failed_retryable
	_ = err
	return state.StatusFailedRetryable
}

func stringPtr(s string) *string {
	return &s
}
