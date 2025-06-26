package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds PostgreSQL connection settings.
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// LoadConfig reads settings from environment variables with fallbacks.
func LoadConfig() Config {
	return Config{
		User:     getenv("POSTGRES_USER", "mem0"),
		Password: getenv("POSTGRES_PASSWORD", "mem0pass"),
		Host:     getenv("POSTGRES_HOST", "localhost"),
		Port:     getenv("POSTGRES_PORT", "5432"),
		Name:     getenv("POSTGRES_DB", "mem0"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ConnString builds a Postgres connection string.
func (c Config) ConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.Name)
}

// Connect opens a pgx connection pool using the provided config.
func Connect(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnString())
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
