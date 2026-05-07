package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Service exited.")

}
