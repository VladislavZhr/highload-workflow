package main

import (
	"context"
	"database/sql"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/VladislavZhr/highload-workflow/handler/internal/config"
	"github.com/VladislavZhr/highload-workflow/handler/internal/kafka"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/finalize"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/mapper"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/start"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/state"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := newDB(cfg)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	stateRepo := state.NewRepository(db)

	startService := start.NewService(stateRepo, cfg.MaxRetryCount, cfg.LeaseDuration)

	domainMapper := mapper.NewMapper()
	mapperService := mapper.NewService(domainMapper)
	finalizeService := finalize.NewService(stateRepo)

	orchestrator := pipeline.NewOrchestrator(
		db,
		startService,
		mapperService,
		finalizeService,
	)

	consumerCfg := kafka.Config{
		Brokers:          cfg.KafkaBrokers,
		GroupID:          cfg.KafkaGroupID,
		Topic:            cfg.KafkaTopic,
		MinBytes:         cfg.KafkaMinBytes,
		MaxBytes:         cfg.KafkaMaxBytes,
		MaxWait:          cfg.KafkaMaxWait,
		WorkersCount:     cfg.WorkersCount,
		JobsBufferSize:   cfg.JobsBufferSize,
		ResultsBufSize:   cfg.ResultsBufferSize,
		ReadBatchTimeout: cfg.ReadBatchTimeout,
		CommitTimeout:    cfg.CommitTimeout,
		ShutdownTimeout:  cfg.ShutdownTimeout,
	}

	consumer, err := kafka.NewConsumer(consumerCfg, orchestrator)
	if err != nil {
		log.Fatalf("init kafka consumer: %v", err)
	}

	log.Printf(
		"handler started kafka_brokers=%v group_id=%s topic=%s workers=%d",
		cfg.KafkaBrokers,
		cfg.KafkaGroupID,
		cfg.KafkaTopic,
		cfg.WorkersCount,
	)

	if err := consumer.Run(ctx); err != nil {
		log.Fatalf("consumer stopped with error: %v", err)
	}

	log.Println("consumer stopped")
}

func newDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(context.Background(), cfg.DBPingTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
