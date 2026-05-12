package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type DeviceEvent struct {
	DeviceID  string  `json:"device_id"`
	Timestamp int64   `json:"timestamp"`
	AccX      float64 `json:"acc_x"`
	AccY      float64 `json:"acc_y"`
	AccZ      float64 `json:"acc_z"`
	PGA       float64 `json:"pga"`
	StaLta    float64 `json:"sta_lta"`
	IsTrigger bool    `json:"is_trigger"`
}

func KafkaReader(ctx context.Context, db *sql.DB) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "user-events",
		StartOffset: kafka.FirstOffset,
	})
	defer func() {
		if err := reader.Close(); err != nil {
			fmt.Printf("Failed to close Kafka reader: %v\n", err)
		}
	}()
	fmt.Println("masuk ke kafka")

	for {
		fmt.Println("ctx err before loop:", ctx.Err())
		fmt.Println("before read message")
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("Kafka Reader stopping via context...")
				return
			}
			fmt.Printf("Kafka Error: %v\n", err)
			return
		}

		var event DeviceEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			fmt.Printf("Failed to unmarshal Kafka message: %v\n", err)
			continue
		}

		fmt.Println("loop kafka")

		recordedAt := time.UnixMicro(event.Timestamp)

		query := `
			INSERT INTO sensor_readings (
				sensor_id, 
				recorded_at, 
				acc_x, 
				acc_y, 
				acc_z, 
				pga, 
				sta_lta, 
				is_trigger
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		fmt.Println("before exec")
		_, err = db.ExecContext(ctx, query,
			event.DeviceID,
			recordedAt,
			event.AccX,
			event.AccY,
			event.AccZ,
			event.PGA,
			event.StaLta,
			event.IsTrigger,
		)
		fmt.Println("after exec")

		if err != nil {
			fmt.Printf("Failed to insert into sensor_readings: %v\n", err)
			continue
		}

		fmt.Printf("Inserted reading from sensor %s at %v\n", event.DeviceID, recordedAt)
	}
}
