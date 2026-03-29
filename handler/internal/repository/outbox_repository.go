package repository

import (
	"context"
	"database/sql"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
)

type OutboxRepository struct {
	db *sql.DB
}

func NewOutboxRepository(db *sql.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) Insert(
	ctx context.Context,
	tx *sql.Tx,
	message model.OutboxMessage,
) error {
	query := `
		INSERT INTO outbox_messages (
			correlation_id,
			request_id
		)
		VALUES ($1, $2)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		message.CorrelationID,
		message.RequestID,
	)
	if err != nil {
		return err
	}

	return nil
}
