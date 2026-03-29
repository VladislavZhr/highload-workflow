package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline"
	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline/start"
)

type Consumer struct {
	reader       *kafkago.Reader
	orchestrator *pipeline.Orchestrator
	cfg          Config

	jobs    chan Job
	results chan Result

	workersWG sync.WaitGroup
	resultsWG sync.WaitGroup
}

func NewConsumer(cfg Config, orchestrator *pipeline.Orchestrator) (*Consumer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: cfg.MinBytes,
		MaxBytes: cfg.MaxBytes,
		MaxWait:  cfg.MaxWait,
	})

	return &Consumer{
		reader:       reader,
		orchestrator: orchestrator,
		cfg:          cfg,
		jobs:         make(chan Job, cfg.JobsBufferSize),
		results:      make(chan Result, cfg.ResultsBufSize),
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	defer func() {
		if err := c.reader.Close(); err != nil {
			log.Printf("kafka reader close error: %v", err)
		}
	}()

	c.startWorkers(ctx)
	c.startResultLoop(ctx)

	readErr := c.readLoop(ctx)

	close(c.jobs)
	c.workersWG.Wait()

	close(c.results)
	c.resultsWG.Wait()

	if readErr != nil && !errors.Is(readErr, context.Canceled) {
		return readErr
	}

	return nil
}

func (c *Consumer) startWorkers(ctx context.Context) {
	for i := 0; i < c.cfg.WorkersCount; i++ {
		c.workersWG.Add(1)

		workerID := i + 1

		go func() {
			defer c.workersWG.Done()
			worker(ctx, workerID, c.orchestrator, c.jobs, c.results)
		}()
	}
}

func (c *Consumer) startResultLoop(ctx context.Context) {
	c.resultsWG.Add(1)

	go func() {
		defer c.resultsWG.Done()
		c.resultLoop(ctx)
	}()
}

func (c *Consumer) readLoop(ctx context.Context) error {
	for {
		msg, err := c.fetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}

			return fmt.Errorf("fetch kafka message: %w", err)
		}

		job := Job{
			Message: msg,
		}

		select {
		case <-ctx.Done():
			return ctx.Err()

		case c.jobs <- job:
			log.Printf(
				"message accepted partition=%d offset=%d",
				msg.Partition,
				msg.Offset,
			)
		}
	}
}

func (c *Consumer) fetchMessage(ctx context.Context) (kafkago.Message, error) {
	readCtx, cancel := context.WithTimeout(ctx, c.cfg.ReadBatchTimeout)
	defer cancel()

	msg, err := c.reader.FetchMessage(readCtx)
	if err != nil {
		return kafkago.Message{}, err
	}

	return msg, nil
}

func (c *Consumer) resultLoop(ctx context.Context) {
	for result := range c.results {
		if result.Err != nil {
			if shouldCommit(result.Err) {
				if err := c.commitMessage(ctx, result.Message); err != nil {
					log.Printf(
						"commit failed after handled error partition=%d offset=%d err=%v",
						result.Message.Partition,
						result.Message.Offset,
						err,
					)
					continue
				}

				log.Printf(
					"message committed after handled error partition=%d offset=%d err=%v",
					result.Message.Partition,
					result.Message.Offset,
					result.Err,
				)
				continue
			}

			log.Printf(
				"message not committed partition=%d offset=%d err=%v",
				result.Message.Partition,
				result.Message.Offset,
				result.Err,
			)
			continue
		}

		if err := c.commitMessage(ctx, result.Message); err != nil {
			log.Printf(
				"commit failed partition=%d offset=%d err=%v",
				result.Message.Partition,
				result.Message.Offset,
				err,
			)
			continue
		}

		log.Printf(
			"message committed partition=%d offset=%d correlation_id=%s",
			result.Message.Partition,
			result.Message.Offset,
			result.Transport.Message.Header.CorrelationID,
		)
	}
}

func (c *Consumer) commitMessage(ctx context.Context, msg kafkago.Message) error {
	commitCtx, cancel := context.WithTimeout(ctx, c.cfg.CommitTimeout)
	defer cancel()

	if err := c.reader.CommitMessages(commitCtx, msg); err != nil {
		return fmt.Errorf("commit kafka message: %w", err)
	}

	return nil
}

func shouldCommit(err error) bool {
	if errors.Is(err, start.ErrSkipCompleted) {
		return true
	}

	if errors.Is(err, start.ErrSkipPermanent) {
		return true
	}

	if isMalformedMessageError(err) {
		return true
	}

	return false
}
