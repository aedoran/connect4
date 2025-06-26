package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines persistence operations used by the app.
type Repository interface {
	CreateUser(ctx context.Context, username string) (int64, error)
	CreateMemory(ctx context.Context, userID int64, content string) (int64, error)
	AddEmbedding(ctx context.Context, memoryID int64, vector []float32) error
}

// PgxRepository implements Repository with a pgx pool.
type PgxRepository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *PgxRepository { return &PgxRepository{pool: pool} }

func (r *PgxRepository) CreateUser(ctx context.Context, username string) (int64, error) {
	row := r.pool.QueryRow(ctx, "INSERT INTO users (username) VALUES ($1) RETURNING id", username)
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PgxRepository) CreateMemory(ctx context.Context, userID int64, content string) (int64, error) {
	row := r.pool.QueryRow(ctx, "INSERT INTO memories (user_id, content) VALUES ($1,$2) RETURNING id", userID, content)
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PgxRepository) AddEmbedding(ctx context.Context, memoryID int64, vector []float32) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO embeddings (memory_id, vector) VALUES ($1,$2)", memoryID, vector)
	return err
}
