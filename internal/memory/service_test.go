package memory

import (
	"context"
	"fmt"
	"testing"

	"mem0-go/internal/graph"
	"mem0-go/internal/vector"
)

type stubRepo struct {
	users      []string
	memories   []string
	embeddings [][]float32
	createErr  error
	embedErr   error
}

func (s *stubRepo) CreateUser(ctx context.Context, username string) (int64, error) {
	s.users = append(s.users, username)
	return int64(len(s.users)), nil
}

func (s *stubRepo) CreateMemory(ctx context.Context, userID int64, content string) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}
	s.memories = append(s.memories, content)
	return int64(len(s.memories)), nil
}

func (s *stubRepo) AddEmbedding(ctx context.Context, memoryID int64, vec []float32) error {
	if s.embedErr != nil {
		return s.embedErr
	}
	s.embeddings = append(s.embeddings, vec)
	return nil
}

type stubVector struct {
	upsertCalled bool
	queryCalled  bool
	upsertErr    error
	queryErr     error
}

func (s *stubVector) Upsert(ctx context.Context, col string, pts []vector.Point) error {
	s.upsertCalled = true
	return s.upsertErr
}

func (s *stubVector) Query(ctx context.Context, col string, vec []float32, limit int) ([]vector.QueryResult, error) {
	s.queryCalled = true
	if s.queryErr != nil {
		return nil, s.queryErr
	}
	return []vector.QueryResult{{ID: "1", Score: 0.9}}, nil
}

type stubGraph struct {
	nodes   map[string]graph.Node
	edges   []graph.Edge
	next    int
	nodeErr error
	edgeErr error
}

func (g *stubGraph) CreateNode(_ context.Context, label string, props map[string]interface{}) (string, error) {
	if g.nodeErr != nil {
		return "", g.nodeErr
	}
	if g.nodes == nil {
		g.nodes = make(map[string]graph.Node)
	}
	g.next++
	id := fmt.Sprintf("n%d", g.next)
	g.nodes[id] = graph.Node{ID: id, Label: label, Props: props}
	return id, nil
}

func (g *stubGraph) CreateEdge(_ context.Context, from, to, relType string, props map[string]interface{}) (string, error) {
	if g.edgeErr != nil {
		return "", g.edgeErr
	}
	g.next++
	id := fmt.Sprintf("e%d", g.next)
	g.edges = append(g.edges, graph.Edge{ID: id, From: from, To: to, Type: relType, Props: props})
	return id, nil
}

func (g *stubGraph) Neighbors(_ context.Context, id, relType string) ([]graph.Node, error) {
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

func TestStoreAndSearch(t *testing.T) {
	repo := &stubRepo{}
	vec := &stubVector{}
	g := &stubGraph{}
	svc := NewService(repo, vec, g)

	id, err := svc.StoreMemory(context.Background(), 1, "hello", []float32{1, 2})
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	if id != 1 {
		t.Fatalf("unexpected id %d", id)
	}
	if !vec.upsertCalled {
		t.Fatalf("upsert not called")
	}

	res, err := svc.Search(context.Background(), []float32{1, 2}, 1)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(res) != 1 || res[0].ID != 1 {
		t.Fatalf("unexpected results: %+v", res)
	}
	if !vec.queryCalled {
		t.Fatalf("query not called")
	}
}

func TestRelateEntities(t *testing.T) {
	repo := &stubRepo{}
	vec := &stubVector{}
	g := &stubGraph{}
	svc := NewService(repo, vec, g)

	n1, _ := svc.CreateEntity(context.Background(), "Person", nil)
	n2, _ := svc.CreateEntity(context.Background(), "Person", nil)
	id, err := svc.RelateEntities(context.Background(), n1, n2, "KNOWS", nil)
	if err != nil {
		t.Fatalf("relate: %v", err)
	}
	if id == "" {
		t.Fatalf("empty edge id")
	}
	neigh, _ := g.Neighbors(context.Background(), n1, "KNOWS")
	if len(neigh) != 1 || neigh[0].ID != n2 {
		t.Fatalf("unexpected neighbors: %+v", neigh)
	}
}

func TestStoreMemoryErrors(t *testing.T) {
	repo := &stubRepo{createErr: fmt.Errorf("boom")}
	vec := &stubVector{}
	g := &stubGraph{}
	if _, err := NewService(repo, vec, g).StoreMemory(context.Background(), 1, "x", nil); err == nil {
		t.Fatalf("expected create error")
	}

	repo.createErr = nil
	repo.embedErr = fmt.Errorf("bad")
	if _, err := NewService(repo, vec, g).StoreMemory(context.Background(), 1, "x", nil); err == nil {
		t.Fatalf("expected embed error")
	}

	repo.embedErr = nil
	vec.upsertErr = fmt.Errorf("up")
	if _, err := NewService(repo, vec, g).StoreMemory(context.Background(), 1, "x", nil); err == nil {
		t.Fatalf("expected upsert error")
	}
}

func TestSearchError(t *testing.T) {
	repo := &stubRepo{}
	vec := &stubVector{queryErr: fmt.Errorf("q")}
	g := &stubGraph{}
	svc := NewService(repo, vec, g)
	if _, err := svc.Search(context.Background(), nil, 1); err == nil {
		t.Fatalf("expected query error")
	}
}
