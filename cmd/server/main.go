package main

import (
	"log"

	"github.com/Wild-sergunys/shrtic/internal/config"
	"github.com/Wild-sergunys/shrtic/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	db, err := database.NewPostgres(&cfg.DB)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(cfg.DB.MigrateDSN()); err != nil {
		log.Fatalf("Ошибка миграций: %v", err)
	}

	redisClient, err := database.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("Все сервисы запущены")
}
