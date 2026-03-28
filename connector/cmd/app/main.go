package main

import (
	"log"
	"net/http"

	"github.com/VladislavZhr/highload-workflow/connector/internal/handler"
	"github.com/VladislavZhr/highload-workflow/connector/internal/kafka"
	"github.com/VladislavZhr/highload-workflow/connector/internal/service"
)

func main() {
	producer := kafka.NewProducer(
		[]string{"localhost:9092"},
		"connector-topic",
	)
	defer producer.Close()

	connectorService := service.NewConnectorService(producer)
	httpHandler := handler.NewHTTPHandler(connectorService)

	http.HandleFunc("/process", httpHandler.HandleProcessRequest)

	log.Println("connector started on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
