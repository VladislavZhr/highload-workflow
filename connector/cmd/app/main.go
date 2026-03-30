package main

import (
	"log"
	"net/http"

	"github.com/VladislavZhr/highload-workflow/connector/internal/config"
	"github.com/VladislavZhr/highload-workflow/connector/internal/handler"
	"github.com/VladislavZhr/highload-workflow/connector/internal/kafka"
	"github.com/VladislavZhr/highload-workflow/connector/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	producer := kafka.NewProducer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
	)
	defer producer.Close()

	connectorService := service.NewConnectorService(producer)
	httpHandler := handler.NewHTTPHandler(connectorService)

	http.HandleFunc("/process", httpHandler.HandleProcessRequest)

	log.Printf(
		"connector started addr=%s kafka_brokers=%v kafka_topic=%s",
		cfg.HTTPAddr(),
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
	)

	if err := http.ListenAndServe(cfg.HTTPAddr(), nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
