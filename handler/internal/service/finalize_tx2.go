package service

import (
	"context"
	"database/sql"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
	"github.com/VladislavZhr/highload-workflow/handler/internal/repository"
)

type FinalizeTx2Service struct {
	processingRepo *repository.ProcessingRepository
	outboxRepo     *repository.OutboxRepository
}

func NewFinalizeTx2Service(
	processingRepo *repository.ProcessingRepository,
	outboxRepo *repository.OutboxRepository,
) *FinalizeTx2Service {
	return &FinalizeTx2Service{
		processingRepo: processingRepo,
		outboxRepo:     outboxRepo,
	}
}

func (s *FinalizeTx2Service) Complete(
	ctx context.Context,
	tx *sql.Tx,
	msg model.ProcessedMessage,
) error {
	if err := s.processingRepo.UpdateStatus(
		ctx,
		tx,
		msg.CorrelationID,
		model.StatusCompleted,
		nil,
	); err != nil {
		return err
	}

	if err := s.outboxRepo.Insert(
		ctx,
		tx,
		model.OutboxMessage{
			CorrelationID: msg.CorrelationID,
			RequestID:     msg.RequestID,
		},
	); err != nil {
		return err
	}

	return nil
}
