package kafka

import "time"

type Config struct {
	Brokers          []string
	Topic            string
	GroupID          string
	WorkersCount     int
	JobsBufferSize   int
	ResultsBufSize   int
	MinBytes         int
	MaxBytes         int
	MaxWait          time.Duration
	ReadBatchTimeout time.Duration
	CommitTimeout    time.Duration
	ShutdownTimeout  time.Duration
}

func DefaultConfig() Config {
	return Config{
		Brokers:          []string{"localhost:9092"},
		Topic:            "handler-topic",
		GroupID:          "handler-consumer-group",
		WorkersCount:     8,
		JobsBufferSize:   32,
		ResultsBufSize:   32,
		MinBytes:         1,
		MaxBytes:         10e6, // 10 MB
		MaxWait:          500 * time.Millisecond,
		ReadBatchTimeout: 3 * time.Second,
		CommitTimeout:    5 * time.Second,
		ShutdownTimeout:  10 * time.Second,
	}
}

func (c Config) Validate() error {
	if len(c.Brokers) == 0 {
		return ErrInvalidConfig("brokers must not be empty")
	}

	if c.Topic == "" {
		return ErrInvalidConfig("topic must not be empty")
	}

	if c.GroupID == "" {
		return ErrInvalidConfig("group id must not be empty")
	}

	if c.WorkersCount <= 0 {
		return ErrInvalidConfig("workers count must be greater than 0")
	}

	if c.JobsBufferSize <= 0 {
		return ErrInvalidConfig("jobs buffer size must be greater than 0")
	}

	if c.ResultsBufSize <= 0 {
		return ErrInvalidConfig("results buffer size must be greater than 0")
	}

	if c.MinBytes <= 0 {
		return ErrInvalidConfig("min bytes must be greater than 0")
	}

	if c.MaxBytes < c.MinBytes {
		return ErrInvalidConfig("max bytes must be greater than or equal to min bytes")
	}

	if c.MaxWait <= 0 {
		return ErrInvalidConfig("max wait must be greater than 0")
	}

	if c.ReadBatchTimeout <= 0 {
		return ErrInvalidConfig("read batch timeout must be greater than 0")
	}

	if c.CommitTimeout <= 0 {
		return ErrInvalidConfig("commit timeout must be greater than 0")
	}

	if c.ShutdownTimeout <= 0 {
		return ErrInvalidConfig("shutdown timeout must be greater than 0")
	}

	return nil
}
