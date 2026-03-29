package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/VladislavZhr/highload-workflow/handler/internal/pipeline"
	"github.com/VladislavZhr/highload-workflow/handler/internal/transport"
)

func worker(ctx context.Context, workerID int, orchestrator *pipeline.Orchestrator, jobs <-chan Job, results chan<- Result) {

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			tm, err := decodeTransportMessage(job.Message.Value)
			if err != nil {
				results <- Result{
					Message: job.Message,
					Err:     err,
				}

				log.Printf(
					"worker=%d partition=%d offset=%d decode_err=%v",
					workerID,
					job.Message.Partition,
					job.Message.Offset,
					err,
				)
				continue
			}

			_, err = orchestrator.Handle(ctx, tm)

			results <- Result{
				Message:   job.Message,
				Transport: tm,
				Err:       err,
			}

			log.Printf(
				"worker=%d partition=%d offset=%d correlation_id=%s err=%v",
				workerID,
				job.Message.Partition,
				job.Message.Offset,
				tm.Message.Header.CorrelationID,
				err,
			)
		}
	}
}

func decodeTransportMessage(payload []byte) (transport.TransportMessage, error) {
	var tm transport.TransportMessage

	if len(payload) == 0 {
		return transport.TransportMessage{}, fmt.Errorf("%w: empty kafka message payload", ErrMalformedMessage)
	}

	if err := json.Unmarshal(payload, &tm); err != nil {
		return transport.TransportMessage{}, fmt.Errorf("%w: decode transport message: %v", ErrMalformedMessage, err)
	}

	if tm.Message.Header.CorrelationID == "" {
		return transport.TransportMessage{}, fmt.Errorf("%w: transport message correlationId is empty", ErrMalformedMessage)
	}

	if tm.Message.Header.RequestID == "" {
		return transport.TransportMessage{}, fmt.Errorf("%w: transport message requestId is empty", ErrMalformedMessage)
	}

	return tm, nil
}

func isMalformedMessageError(err error) bool {
	return errors.Is(err, ErrMalformedMessage)
}
