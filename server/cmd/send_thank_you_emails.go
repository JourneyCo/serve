package main

import (
	"context"
	"flag"
	"log"

	"github.com/joho/godotenv"
	"serve/config"
	"serve/database"
	"serve/services"
)

func main() {
	// Parse command line flags
	var envFile string
	flag.StringVar(&envFile, "env", ".env", "Path to environment file")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: Could not load %s file: %v", envFile, err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize email service
	emailService := services.NewEmailService(cfg)

	// Send thank you emails
	log.Println("Starting thank you email sending process...")
	ctx := context.Background()
	if err := emailService.SendThankYouToAllUsers(ctx, db); err != nil {
		log.Fatalf("Failed to send thank you emails: %v", err)
	}

	log.Println("Thank you email sending process completed successfully!")
}
