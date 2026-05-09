package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"faiqmain.com/internal"
	"faiqmain.com/internal/database"
	"faiqmain.com/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// db postgres setup
	connStr := os.Getenv("DATABASE_URL_CLIENT")
	db, err := database.NewPostgresConnection(connStr)
	if err != nil {
		fmt.Printf("Database initialization failed: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("Failed to close database: %v\n", err)
		}
	}()

	// kafka reader
	go internal.KafkaReader()

	// http server
	mux := server.Setup(db)
	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
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
