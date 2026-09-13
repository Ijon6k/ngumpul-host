package database

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"ngumpul-host/backend/internal/config"
)

//go:embed schema.sql
var schemaSQL string

func ConnectAndMigrate(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database url: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	var pool *pgxpool.Pool
	// Retry connection loop for docker compose readiness
	for i := 0; i < 15; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			pingErr := pool.Ping(ctx)
			if pingErr == nil {
				break
			}
			pool.Close()
			err = pingErr
		}
		log.Printf("Waiting for database connection (attempt %d/15)...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to database after retries: %w", err)
	}

	log.Println("Database connection established. Applying migrations...")

	// Execute migration
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Database migrations applied successfully.")


	return pool, nil
}
