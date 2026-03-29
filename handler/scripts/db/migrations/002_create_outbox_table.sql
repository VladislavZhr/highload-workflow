CREATE TABLE outbox_messages (
                                 id BIGSERIAL PRIMARY KEY,

                                 correlation_id UUID NOT NULL,
                                 request_id UUID NOT NULL,

                                 status TEXT NOT NULL DEFAULT 'pending',

                                 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                 published_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_messages_status_created_at
    ON outbox_messages (status, created_at);