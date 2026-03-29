package start

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertProcessing(ctx context.Context, tx *sql.Tx, record MessageState) (bool, error) {

	query := `
		INSERT INTO message_processing_state (
			correlation_id,
			request_id,
			status,
			lease_until
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (correlation_id) DO NOTHING
	`

	res, err := tx.ExecContext(ctx, query, record.CorrelationID, record.RequestID, record.Status, record.LeaseUntil)

	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *Repository) GetByCorrelationID(ctx context.Context, tx *sql.Tx, correlationID string) (*MessageState, error) {

	query := `
		SELECT
			correlation_id,
			request_id,
			status,
			lease_until,
			retry_count,
			last_error,
			created_at,
			updated_at
		FROM message_processing_state
		WHERE correlation_id = $1
	`

	row := tx.QueryRowContext(ctx, query, correlationID)

	rec := &MessageState{}

	err := row.Scan(
		&rec.CorrelationID,
		&rec.RequestID,
		&rec.Status,
		&rec.LeaseUntil,
		&rec.RetryCount,
		&rec.LastError,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rec, nil
}

func (r *Repository) TryUpdateExpiredLease(ctx context.Context, tx *sql.Tx, correlationID string, leaseUntil any) (bool, error) {
	query := `
		UPDATE message_processing_state
		SET
			lease_until = $2,
			retry_count = retry_count + 1
		WHERE correlation_id = $1
		  AND status = 'processing'
		  AND lease_until IS NOT NULL
		  AND lease_until < NOW()
	`

	res, err := tx.ExecContext(
		ctx,
		query,
		correlationID,
		leaseUntil,
	)

	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, tx *sql.Tx, correlationID string, status ProcessingStatus, lastError *string) error {

	query := `
		UPDATE message_processing_state
		SET
			status = $2,
			last_error = $3,
			lease_until = NULL
		WHERE correlation_id = $1
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		correlationID,
		status,
		lastError,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) TryRetryFailed(ctx context.Context, tx *sql.Tx, correlationID string, maxRetryCount int) (bool, error) {
	query := `
		UPDATE message_processing_state
		SET
			status = 'processing',
			retry_count = retry_count + 1,
			last_error = NULL
		WHERE correlation_id = $1
		  AND status = 'failed_retryable'
		  AND retry_count < $2
	`

	res, err := tx.ExecContext(ctx, query, correlationID, maxRetryCount)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *Repository) MarkFailedPermanentIfRetryExceeded(ctx context.Context, tx *sql.Tx, correlationID string, maxRetryCount int, lastError *string) (bool, error) {
	query := `
		UPDATE message_processing_state
		SET
			status = 'failed_permanent',
			last_error = $3,
			lease_until = NULL
		WHERE correlation_id = $1
		  AND retry_count >= $2
		  AND status IN ('processing', 'failed_retryable')
	`

	res, err := tx.ExecContext(
		ctx,
		query,
		correlationID,
		maxRetryCount,
		lastError,
	)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}
