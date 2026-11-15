package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"hemingway-kayak-web/internal/api"
	"hemingway-kayak-web/internal/services"
	"hemingway-kayak-web/pkg/config"
	"hemingway-kayak-web/pkg/database"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	if err := database.EnsureDatabase(); err != nil {
		log.Fatalf("Failed to ensure database: %v", err)
	}

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := database.SeedData(); err != nil {
		log.Printf("Seed data not loaded: %v", err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redisClient := database.ConnectRedis()

	userService := services.NewUserService(db)
	trackingService := services.NewTrackingService(db, redisClient)
	tournamentService := services.NewTournamentService(db)
	tripService := services.NewTripService(db)
	historyService := services.NewHistoryService(db)

	server := api.NewServer(db, redisClient, userService, trackingService, tournamentService, tripService, historyService)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := server.Shutdown(); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
}
