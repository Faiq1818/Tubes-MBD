package internal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func KafkaReader(ctx context.Context, db *sql.DB) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-events",
		GroupID: "my-service-group",
	})
	defer func() {
		if err := reader.Close(); err != nil {
			fmt.Printf("Failed to close Kafka reader: %v\n", err)
		}
	}()

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("Kafka Reader stopping via context...")
				return
			}
			fmt.Printf("Kafka Error: %v\n", err)
			return
		}

		// insert data from kafka to db
		query := `INSERT INTO user_logs (data, created_at) VALUES ($1, NOW())`

		_, err = db.ExecContext(ctx, query, string(m.Value))
		if err != nil {
			fmt.Printf("Failed to insert Kafka message to DB: %v\n", err)
			continue
		}

		fmt.Printf("Kafka received: %s\n", string(m.Value))
	}
}
