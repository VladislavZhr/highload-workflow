package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDSN string

	KafkaBrokers []string
	KafkaGroupID string
	KafkaTopic   string

	KafkaMinBytes int
	KafkaMaxBytes int
	KafkaMaxWait  time.Duration

	WorkersCount      int
	JobsBufferSize    int
	ResultsBufferSize int

	ReadBatchTimeout time.Duration
	CommitTimeout    time.Duration
	ShutdownTimeout  time.Duration

	MaxRetryCount int
	LeaseDuration time.Duration

	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration
	DBPingTimeout     time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DatabaseDSN: "postgres://postgres:postgres@postgres:5432/highload_workflow?sslmode=disable",

		KafkaBrokers: []string{
			"kafka-1:29092",
			"kafka-2:29092",
			"kafka-3:29092",
		},
		KafkaGroupID: "highload-workflow-handler-group",
		KafkaTopic:   "input-topic",

		KafkaMinBytes: 1,
		KafkaMaxBytes: 100 * 1024 * 1024,
		KafkaMaxWait:  2 * time.Second,

		WorkersCount:      4,
		JobsBufferSize:    64,
		ResultsBufferSize: 64,

		ReadBatchTimeout: 5 * time.Second,
		CommitTimeout:    5 * time.Second,
		ShutdownTimeout:  10 * time.Second,

		MaxRetryCount: 3,
		LeaseDuration: 30 * time.Second,

		DBMaxOpenConns:    20,
		DBMaxIdleConns:    10,
		DBConnMaxLifetime: 30 * time.Minute,
		DBConnMaxIdleTime: 10 * time.Minute,
		DBPingTimeout:     5 * time.Second,
	}

	var err error

	cfg.DatabaseDSN = getEnvString("DATABASE_DSN", cfg.DatabaseDSN)
	cfg.KafkaBrokers = getEnvStringSlice("KAFKA_BROKERS", cfg.KafkaBrokers)
	cfg.KafkaGroupID = getEnvString("KAFKA_GROUP_ID", cfg.KafkaGroupID)
	cfg.KafkaTopic = getEnvString("KAFKA_TOPIC", cfg.KafkaTopic)

	if cfg.KafkaMinBytes, err = getEnvInt("KAFKA_MIN_BYTES", cfg.KafkaMinBytes); err != nil {
		return Config{}, fmt.Errorf("parse KAFKA_MIN_BYTES: %w", err)
	}
	if cfg.KafkaMaxBytes, err = getEnvInt("KAFKA_MAX_BYTES", cfg.KafkaMaxBytes); err != nil {
		return Config{}, fmt.Errorf("parse KAFKA_MAX_BYTES: %w", err)
	}
	if cfg.KafkaMaxWait, err = getEnvDuration("KAFKA_MAX_WAIT", cfg.KafkaMaxWait); err != nil {
		return Config{}, fmt.Errorf("parse KAFKA_MAX_WAIT: %w", err)
	}

	if cfg.WorkersCount, err = getEnvInt("WORKERS_COUNT", cfg.WorkersCount); err != nil {
		return Config{}, fmt.Errorf("parse WORKERS_COUNT: %w", err)
	}
	if cfg.JobsBufferSize, err = getEnvInt("JOBS_BUFFER_SIZE", cfg.JobsBufferSize); err != nil {
		return Config{}, fmt.Errorf("parse JOBS_BUFFER_SIZE: %w", err)
	}
	if cfg.ResultsBufferSize, err = getEnvInt("RESULTS_BUFFER_SIZE", cfg.ResultsBufferSize); err != nil {
		return Config{}, fmt.Errorf("parse RESULTS_BUFFER_SIZE: %w", err)
	}

	if cfg.ReadBatchTimeout, err = getEnvDuration("READ_BATCH_TIMEOUT", cfg.ReadBatchTimeout); err != nil {
		return Config{}, fmt.Errorf("parse READ_BATCH_TIMEOUT: %w", err)
	}
	if cfg.CommitTimeout, err = getEnvDuration("COMMIT_TIMEOUT", cfg.CommitTimeout); err != nil {
		return Config{}, fmt.Errorf("parse COMMIT_TIMEOUT: %w", err)
	}
	if cfg.ShutdownTimeout, err = getEnvDuration("SHUTDOWN_TIMEOUT", cfg.ShutdownTimeout); err != nil {
		return Config{}, fmt.Errorf("parse SHUTDOWN_TIMEOUT: %w", err)
	}

	if cfg.MaxRetryCount, err = getEnvInt("MAX_RETRY_COUNT", cfg.MaxRetryCount); err != nil {
		return Config{}, fmt.Errorf("parse MAX_RETRY_COUNT: %w", err)
	}
	if cfg.LeaseDuration, err = getEnvDuration("LEASE_DURATION", cfg.LeaseDuration); err != nil {
		return Config{}, fmt.Errorf("parse LEASE_DURATION: %w", err)
	}

	if cfg.DBMaxOpenConns, err = getEnvInt("DB_MAX_OPEN_CONNS", cfg.DBMaxOpenConns); err != nil {
		return Config{}, fmt.Errorf("parse DB_MAX_OPEN_CONNS: %w", err)
	}
	if cfg.DBMaxIdleConns, err = getEnvInt("DB_MAX_IDLE_CONNS", cfg.DBMaxIdleConns); err != nil {
		return Config{}, fmt.Errorf("parse DB_MAX_IDLE_CONNS: %w", err)
	}
	if cfg.DBConnMaxLifetime, err = getEnvDuration("DB_CONN_MAX_LIFETIME", cfg.DBConnMaxLifetime); err != nil {
		return Config{}, fmt.Errorf("parse DB_CONN_MAX_LIFETIME: %w", err)
	}
	if cfg.DBConnMaxIdleTime, err = getEnvDuration("DB_CONN_MAX_IDLE_TIME", cfg.DBConnMaxIdleTime); err != nil {
		return Config{}, fmt.Errorf("parse DB_CONN_MAX_IDLE_TIME: %w", err)
	}
	if cfg.DBPingTimeout, err = getEnvDuration("DB_PING_TIMEOUT", cfg.DBPingTimeout); err != nil {
		return Config{}, fmt.Errorf("parse DB_PING_TIMEOUT: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.DatabaseDSN) == "" {
		return errors.New("DATABASE_DSN must not be empty")
	}
	if len(c.KafkaBrokers) == 0 {
		return errors.New("KAFKA_BROKERS must not be empty")
	}
	if strings.TrimSpace(c.KafkaGroupID) == "" {
		return errors.New("KAFKA_GROUP_ID must not be empty")
	}
	if strings.TrimSpace(c.KafkaTopic) == "" {
		return errors.New("KAFKA_TOPIC must not be empty")
	}
	if c.KafkaMinBytes <= 0 {
		return errors.New("KAFKA_MIN_BYTES must be greater than 0")
	}
	if c.KafkaMaxBytes < c.KafkaMinBytes {
		return errors.New("KAFKA_MAX_BYTES must be greater than or equal to KAFKA_MIN_BYTES")
	}
	if c.KafkaMaxWait <= 0 {
		return errors.New("KAFKA_MAX_WAIT must be greater than 0")
	}
	if c.WorkersCount <= 0 {
		return errors.New("WORKERS_COUNT must be greater than 0")
	}
	if c.JobsBufferSize <= 0 {
		return errors.New("JOBS_BUFFER_SIZE must be greater than 0")
	}
	if c.ResultsBufferSize <= 0 {
		return errors.New("RESULTS_BUFFER_SIZE must be greater than 0")
	}
	if c.ReadBatchTimeout <= 0 {
		return errors.New("READ_BATCH_TIMEOUT must be greater than 0")
	}
	if c.CommitTimeout <= 0 {
		return errors.New("COMMIT_TIMEOUT must be greater than 0")
	}
	if c.ShutdownTimeout <= 0 {
		return errors.New("SHUTDOWN_TIMEOUT must be greater than 0")
	}
	if c.MaxRetryCount < 0 {
		return errors.New("MAX_RETRY_COUNT must be greater than or equal to 0")
	}
	if c.LeaseDuration <= 0 {
		return errors.New("LEASE_DURATION must be greater than 0")
	}
	if c.DBMaxOpenConns <= 0 {
		return errors.New("DB_MAX_OPEN_CONNS must be greater than 0")
	}
	if c.DBMaxIdleConns < 0 {
		return errors.New("DB_MAX_IDLE_CONNS must be greater than or equal to 0")
	}
	if c.DBConnMaxLifetime <= 0 {
		return errors.New("DB_CONN_MAX_LIFETIME must be greater than 0")
	}
	if c.DBConnMaxIdleTime <= 0 {
		return errors.New("DB_CONN_MAX_IDLE_TIME must be greater than 0")
	}
	if c.DBPingTimeout <= 0 {
		return errors.New("DB_PING_TIMEOUT must be greater than 0")
	}

	return nil
}

func getEnvString(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getEnvStringSlice(key string, fallback []string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}

	parts := strings.Split(v, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	if len(result) == 0 {
		return fallback
	}

	return result
}

func getEnvInt(key string, fallback int) (int, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback, nil
	}
	return strconv.Atoi(v)
}

func getEnvDuration(key string, fallback time.Duration) (time.Duration, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback, nil
	}
	return time.ParseDuration(v)
}
