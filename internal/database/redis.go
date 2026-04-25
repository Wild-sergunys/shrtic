package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/Wild-sergunys/shrtic/internal/config"
)

func NewRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:   0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Println("Redis connected")
	return client, nil
}
