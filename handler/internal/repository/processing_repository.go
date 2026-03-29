package repository

import (
	"context"
	"database/sql"

	"github.com/VladislavZhr/highload-workflow/handler/internal/model"
)

type ProcessingRepository struct {
	db *sql.DB
}

func NewProcessingRepository(db *sql.DB) *ProcessingRepository {
	return &ProcessingRepository{db: db}
}

// InsertProcessing tries to insert new processing record.
// Returns true if inserted, false if conflict (already exists).
func (r *ProcessingRepository) InsertProcessing(
	ctx context.Context,
	tx *sql.Tx,
	record model.ProcessingRecord,
) (bool, error) {

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

	res, err := tx.ExecContext(
		ctx,
		query,
		record.CorrelationID,
		record.RequestID,
		record.Status,
		record.LeaseUntil,
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

// GetByCorrelationID returns processing record by correlation_id.
func (r *ProcessingRepository) GetByCorrelationID(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
) (*model.ProcessingRecord, error) {

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

	var rec model.ProcessingRecord

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

	return &rec, nil
}

func (r *ProcessingRepository) AcquireExpiredProcessing(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
	requestID string,
	leaseUntil any,
) (bool, error) {

	query := `
		UPDATE message_processing_state
		SET
			request_id = $2,
			status = 'processing',
			lease_until = $3,
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
		requestID,
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

func (r *ProcessingRepository) UpdateStatus(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
	status model.ProcessingStatus,
	lastError *string,
) error {

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

func (r *ProcessingRepository) AcquireRetryableFailed(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
	requestID string,
	leaseUntil any,
) (bool, error) {
	query := `
		UPDATE message_processing_state
		SET
			request_id = $2,
			status = 'processing',
			lease_until = $3,
			retry_count = retry_count + 1,
			last_error = NULL
		WHERE correlation_id = $1
		  AND status = 'failed_retryable'
	`

	res, err := tx.ExecContext(
		ctx,
		query,
		correlationID,
		requestID,
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

func (r *ProcessingRepository) MarkFailedPermanentIfRetryExceeded(
	ctx context.Context,
	tx *sql.Tx,
	correlationID string,
	maxRetryCount int,
	lastError *string,
) (bool, error) {
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
