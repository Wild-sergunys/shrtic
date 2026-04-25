package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Wild-sergunys/shrtic/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(cfg *config.DBConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	log.Println("PostgreSQL connected")
	return pool, nil
}
