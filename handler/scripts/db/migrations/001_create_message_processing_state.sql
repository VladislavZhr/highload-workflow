CREATE TABLE message_processing_state (
  correlation_id      UUID PRIMARY KEY,
  request_id          UUID NOT NULL,
  status              VARCHAR(32) NOT NULL,
  lease_until         TIMESTAMPTZ NULL,
  retry_count         INTEGER NOT NULL DEFAULT 0,
  last_error          TEXT NULL,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_message_processing_state_status
      CHECK (status IN (
                        'processing',
                        'completed',
                        'failed_permanent',
                        'failed_retryable'
          )),

  CONSTRAINT chk_message_processing_state_retry_count
      CHECK (retry_count >= 0)
);

CREATE INDEX idx_message_processing_state_status
    ON message_processing_state (status);

CREATE INDEX idx_message_processing_state_lease_until
    ON message_processing_state (lease_until);

CREATE INDEX idx_message_processing_state_updated_at
    ON message_processing_state (updated_at);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_message_processing_state_set_updated_at
    BEFORE UPDATE ON message_processing_state
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();