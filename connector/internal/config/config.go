package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort     int
	KafkaBrokers []string
	KafkaTopic   string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		HTTPPort:     8080,
		KafkaBrokers: []string{"kafka-1:29092", "kafka-2:29092", "kafka-3:29092"},
		KafkaTopic:   "input-topic",
	}

	if v := strings.TrimSpace(os.Getenv("HTTP_PORT")); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("parse HTTP_PORT: %w", err)
		}
		cfg.HTTPPort = port
	}

	if v := strings.TrimSpace(os.Getenv("KAFKA_BROKERS")); v != "" {
		cfg.KafkaBrokers = splitAndTrim(v)
	}

	if v := strings.TrimSpace(os.Getenv("KAFKA_TOPIC")); v != "" {
		cfg.KafkaTopic = v
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.HTTPPort <= 0 {
		return errors.New("HTTP_PORT must be greater than 0")
	}

	if len(c.KafkaBrokers) == 0 {
		return errors.New("KAFKA_BROKERS must not be empty")
	}

	for _, broker := range c.KafkaBrokers {
		if strings.TrimSpace(broker) == "" {
			return errors.New("KAFKA_BROKERS contains empty broker value")
		}
	}

	if strings.TrimSpace(c.KafkaTopic) == "" {
		return errors.New("KAFKA_TOPIC must not be empty")
	}

	return nil
}

func (c Config) HTTPAddr() string {
	return fmt.Sprintf(":%d", c.HTTPPort)
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
