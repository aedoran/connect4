package memory

import (
	"context"
	"fmt"
	"strconv"

	"mem0-go/internal/db"
	"mem0-go/internal/graph"
	"mem0-go/internal/vector"
)

// Service orchestrates storage, search and relationships across
// Postgres, Qdrant and Neo4j.
// vectorStore defines the subset of vector.Client used by Service.
type vectorStore interface {
	Upsert(ctx context.Context, collection string, pts []vector.Point) error
	Query(ctx context.Context, collection string, vector []float32, limit int) ([]vector.QueryResult, error)
}

type graphStore interface {
	CreateNode(ctx context.Context, label string, props map[string]interface{}) (string, error)
	CreateEdge(ctx context.Context, from, to, relType string, props map[string]interface{}) (string, error)
	Neighbors(ctx context.Context, id, relType string) ([]graph.Node, error)
}

type Service struct {
	repo   db.Repository
	vector vectorStore
	graph  graphStore
}

// NewService constructs a Service.
func NewService(repo db.Repository, v vectorStore, g graphStore) *Service {
	return &Service{repo: repo, vector: v, graph: g}
}

// StoreMemory persists the text and embedding then indexes it in Qdrant.
func (s *Service) StoreMemory(ctx context.Context, userID int64, content string, emb []float32) (int64, error) {
	id, err := s.repo.CreateMemory(ctx, userID, content)
	if err != nil {
		return 0, err
	}
	if err := s.repo.AddEmbedding(ctx, id, emb); err != nil {
		return 0, err
	}
	if err := s.vector.Upsert(ctx, "memories", []vector.Point{{ID: fmt.Sprint(id), Vector: emb}}); err != nil {
		return 0, err
	}
	return id, nil
}

// MemoryResult represents a search match.
type MemoryResult struct {
	ID    int64
	Score float32
}

// Search returns similar memories using Qdrant.
func (s *Service) Search(ctx context.Context, emb []float32, limit int) ([]MemoryResult, error) {
	res, err := s.vector.Query(ctx, "memories", emb, limit)
	if err != nil {
		return nil, err
	}
	out := make([]MemoryResult, 0, len(res))
	for _, r := range res {
		id, err := strconv.ParseInt(r.ID, 10, 64)
		if err != nil {
			continue
		}
		out = append(out, MemoryResult{ID: id, Score: r.Score})
	}
	return out, nil
}

// CreateEntity inserts a node into the graph.
func (s *Service) CreateEntity(ctx context.Context, label string, props map[string]interface{}) (string, error) {
	return s.graph.CreateNode(ctx, label, props)
}

// RelateEntities creates a relationship between two nodes.
func (s *Service) RelateEntities(ctx context.Context, fromID, toID, relType string, props map[string]interface{}) (string, error) {
	return s.graph.CreateEdge(ctx, fromID, toID, relType, props)
}
