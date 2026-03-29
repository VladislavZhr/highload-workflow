package finalize

import (
	"context"
	"database/sql"

	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/state"
)

type Service struct {
	repo *state.Repository
}

func NewService(repo *state.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FinalizeFailedState(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
	status state.ProcessingStatus,
	lastError *string,
) error {
	return s.repo.UpdateStatus(
		ctx,
		tx,
		correlationID,
		status,
		lastError,
	)
}

func (s *Service) FinalizeSuccessState(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
) error {
	return s.repo.UpdateStatus(
		ctx,
		tx,
		correlationID,
		state.StatusCompleted,
		nil,
	)
}
