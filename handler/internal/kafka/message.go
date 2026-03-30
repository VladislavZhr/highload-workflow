package kafka

import (
	kafkago "github.com/segmentio/kafka-go"

	"github.com/VladislavZhr/highload-workflow/handler/internal/transport"
)

type Job struct {
	Message kafkago.Message
}

type Result struct {
	Message   kafkago.Message
	Transport transport.TransportMessage
	Err       error
}
