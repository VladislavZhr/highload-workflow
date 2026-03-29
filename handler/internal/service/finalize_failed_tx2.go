package service

import (
	"context"
	"database/sql"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
	"github.com/VladislavZhr/highload-workflow/handler/internal/repository"
)

type FinalizeFailedTx2Service struct {
	processingRepo *repository.ProcessingRepository
}

func NewFinalizeFailedTx2Service(
	processingRepo *repository.ProcessingRepository,
) *FinalizeFailedTx2Service {
	return &FinalizeFailedTx2Service{
		processingRepo: processingRepo,
	}
}

func (s *FinalizeFailedTx2Service) Finalize(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
	status model.ProcessingStatus,
	lastError *string,
) error {
	return s.processingRepo.UpdateStatus(
		ctx,
		tx,
		correlationID,
		status,
		lastError,
	)
}
