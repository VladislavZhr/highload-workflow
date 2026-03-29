package service

import (
	"context"
	"database/sql"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
	"github.com/VladislavZhr/highload-workflow/handler/internal/repository"
)

type FinalizeProcessingService struct {
	processingRepo *repository.ProcessingRepository
}

func NewFinalizeProcessingService(
	processingRepo *repository.ProcessingRepository,
) *FinalizeProcessingService {
	return &FinalizeProcessingService{
		processingRepo: processingRepo,
	}
}

func (s *FinalizeProcessingService) Finalize(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
	result model.ProcessingResult,
) error {
	return s.processingRepo.UpdateStatus(
		ctx,
		tx,
		correlationID,
		result.Status,
		result.LastError,
	)
}
