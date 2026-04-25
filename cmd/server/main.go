package main

import (
	"log"

	"github.com/Wild-sergunys/shrtic/internal/config"
	"github.com/Wild-sergunys/shrtic/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgres(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	redisClient, err := database.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("All services connected successfully")
}
