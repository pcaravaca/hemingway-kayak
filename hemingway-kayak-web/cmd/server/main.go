package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"hemingway-kayak-web/internal/api"
	"hemingway-kayak-web/internal/services"
	"hemingway-kayak-web/pkg/database"
)

func main() {
	if err := database.EnsureDatabase(); err != nil {
		log.Fatalf("Failed to ensure database: %v", err)
	}

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := database.SeedData(); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
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

	server := api.NewServer(redisClient, userService, trackingService, tournamentService, tripService, historyService)

	go server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
