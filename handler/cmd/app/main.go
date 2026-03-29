package main

import (
	"context"
	"database/sql"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

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

	db, err := newDB()
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	stateRepo := state.NewRepository(db)

	// Замінити на свої актуальні значення
	const (
		maxRetryCount = 3
		leaseDuration = 30 * time.Second
	)

	startService := start.NewService(stateRepo, maxRetryCount, leaseDuration)

	// Якщо у тебе конструктор Mapper інший, просто підстав свій.
	// Наприклад:
	// domainMapper := mapper.NewMapper(...)
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
		Brokers: []string{
			"localhost:9092",
		},
		GroupID: "highload-workflow-handler-group",
		Topic:   "input-topic",

		MinBytes: 1,
		MaxBytes: 100 * 1024 * 1024, // 100 MB, бо ти граєшся з payload 70 MB
		MaxWait:  2 * time.Second,

		WorkersCount:   4,
		JobsBufferSize: 64,
		ResultsBufSize: 64,

		ReadBatchTimeout: 5 * time.Second,
		CommitTimeout:    5 * time.Second,
	}

	consumer, err := kafka.NewConsumer(consumerCfg, orchestrator)
	if err != nil {
		log.Fatalf("init kafka consumer: %v", err)
	}

	log.Println("consumer started")

	if err := consumer.Run(ctx); err != nil {
		log.Fatalf("consumer stopped with error: %v", err)
	}

	log.Println("consumer stopped")
}

func newDB() (*sql.DB, error) {
	// Тестовий DSN. Перепишеш під свої параметри.
	dsn := "postgres://postgres:postgres@localhost:5432/highload_workflow?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
