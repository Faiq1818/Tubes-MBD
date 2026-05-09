package internal

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func KafkaReader() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-events",
		GroupID: "my-service-group",
	})
	defer reader.Close()

	for {

		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Kafka Error: %v\n", err)
			return
		}
		fmt.Printf("Kafka received: %s\n", string(m.Value))

	}
}
