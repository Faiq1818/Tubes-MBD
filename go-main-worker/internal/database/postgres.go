package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func NewPostgresConnection(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// defer func() {
	// 	err := db.Close()
	// 	if err != nil {
	// 		fmt.Println("Failed to close database", "error", err)
	// 	}
	// }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	fmt.Println("Successfully setup database")

	return db, nil
}
