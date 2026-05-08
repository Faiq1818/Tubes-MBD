package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func main() {
	err := godotenv.Load()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// db postgres setup
	connStr := os.Getenv("DATABASE_URL_CLIENT")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Database open failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("Failed to close database", "error", err)
		}
	}()

	err = db.PingContext(context.Background())
	if err != nil {
		fmt.Println("Database ping failed", "error", err)
		os.Exit(1)
	}

	// kafka reader
	go func() {
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
	}()

	// http server
	server := &http.Server{Addr: ":8000"}
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go func() {
		fmt.Println("HTTP Server starting on :8000...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP Server Error: %v\n", err)
		}
	}()

	<-sigChan
	fmt.Println("\nShutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Service exited.")
}
