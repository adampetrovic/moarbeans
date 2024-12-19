package main

import (
	"log"
	"github.com/adampetrovic/moarbeans/internal/api"
	"github.com/adampetrovic/moarbeans/internal/auth"
	"github.com/adampetrovic/moarbeans/internal/database"
)

func main() {
	// Initialize database
	db, err := database.NewSQLiteDB("moarbeans.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize email poller
	poller := auth.NewEmailPoller(db)
	go poller.Start()

	// Initialize and start API server
	server := api.NewServer(db)
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
} 