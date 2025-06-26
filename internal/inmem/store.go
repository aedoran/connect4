package inmem

import (
	"context"
	"fmt"

	"mem0-go/internal/db"
	"mem0-go/internal/graph"
	"mem0-go/internal/vector"
)

// Repo implements db.Repository using memory.
// Ensure it satisfies the interface.
var _ db.Repository = (*Repo)(nil)

type Repo struct {
	users      []string
	memories   []string
	embeddings [][]float32
}

func NewRepo() *Repo { return &Repo{} }

func (r *Repo) CreateUser(ctx context.Context, username string) (int64, error) {
	r.users = append(r.users, username)
	return int64(len(r.users)), nil
}

func (r *Repo) CreateMemory(ctx context.Context, userID int64, content string) (int64, error) {
	r.memories = append(r.memories, content)
	return int64(len(r.memories)), nil
}

func (r *Repo) AddEmbedding(ctx context.Context, memoryID int64, vec []float32) error {
	r.embeddings = append(r.embeddings, vec)
	return nil
}

func (r *Repo) GetMemory(ctx context.Context, id int64) (db.Memory, error) {
	if int(id) <= 0 || int(id) > len(r.memories) {
		return db.Memory{}, fmt.Errorf("not found")
	}
	return db.Memory{ID: id, UserID: 1, Content: r.memories[id-1]}, nil
}

// Vector implements vectorStore using memory.
type Vector struct{ points []vector.Point }

func NewVector() *Vector { return &Vector{} }

func (v *Vector) Upsert(ctx context.Context, collection string, pts []vector.Point) error {
	v.points = append(v.points, pts...)
	return nil
}

func (v *Vector) Query(ctx context.Context, collection string, vec []float32, limit int) ([]vector.QueryResult, error) {
	out := []vector.QueryResult{}
	for i, p := range v.points {
		if i >= limit {
			break
		}
		out = append(out, vector.QueryResult{ID: p.ID, Score: 1})
	}
	return out, nil
}

// Graph implements graphStore using in-memory structures.
// We reuse graph.Node and graph.Edge types.
type Graph struct {
	nodes map[string]graph.Node
	edges []graph.Edge
	next  int
}

func NewGraph() *Graph { return &Graph{nodes: make(map[string]graph.Node)} }

func (g *Graph) CreateNode(_ context.Context, label string, props map[string]interface{}) (string, error) {
	g.next++
	id := fmt.Sprintf("n%d", g.next)
	g.nodes[id] = graph.Node{ID: id, Label: label, Props: props}
	return id, nil
}

func (g *Graph) CreateEdge(_ context.Context, from, to, relType string, props map[string]interface{}) (string, error) {
	g.next++
	id := fmt.Sprintf("e%d", g.next)
	g.edges = append(g.edges, graph.Edge{ID: id, From: from, To: to, Type: relType, Props: props})
	return id, nil
}

func (g *Graph) Neighbors(_ context.Context, id, relType string) ([]graph.Node, error) {
	var out []graph.Node
	for _, e := range g.edges {
		if e.Type == relType && e.From == id {
			if n, ok := g.nodes[e.To]; ok {
				out = append(out, n)
			}
		}
	}
	return out, nil
}
