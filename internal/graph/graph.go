package graph

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
)

// Config holds Neo4j connection settings.
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
}

// LoadConfig reads settings from environment variables with fallbacks.
func LoadConfig() Config {
	user := os.Getenv("NEO4J_USER")
	if user == "" {
		user = "neo4j"
	}
	pass := os.Getenv("NEO4J_PASSWORD")
	if pass == "" {
		pass = "neo4jtest"
	}
	host := os.Getenv("NEO4J_HOST")
	if host == "" {
		host = "neo4j"
	}
	port := os.Getenv("NEO4J_PORT")
	if port == "" {
		port = "7474"
	}
	return Config{User: user, Password: pass, Host: host, Port: port}
}

var dialContext = (&net.Dialer{}).DialContext

// Graph provides helpers for storing nodes and edges.
type Graph struct {
	addr  string
	mu    sync.RWMutex
	nodes map[string]Node
	edges []Edge
	next  int
}

// Node represents a graph node.
type Node struct {
	ID    string
	Label string
	Props map[string]interface{}
}

// Edge represents a relationship between two nodes.
type Edge struct {
	ID    string
	From  string
	To    string
	Type  string
	Props map[string]interface{}
}

// Connect verifies a bolt connection can be established and returns a Graph.
func Connect(ctx context.Context, cfg Config) (*Graph, error) {
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	conn, err := dialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	_ = conn.Close()
	return &Graph{
		addr:  addr,
		nodes: make(map[string]Node),
	}, nil
}

// CreateNode inserts a node and returns its generated ID.
func (g *Graph) CreateNode(_ context.Context, label string, props map[string]interface{}) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	id := fmt.Sprintf("n%d", g.next)
	g.nodes[id] = Node{ID: id, Label: label, Props: props}
	return id, nil
}

// CreateEdge inserts a relationship and returns its generated ID.
func (g *Graph) CreateEdge(_ context.Context, from, to, relType string, props map[string]interface{}) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.next++
	id := fmt.Sprintf("e%d", g.next)
	g.edges = append(g.edges, Edge{ID: id, From: from, To: to, Type: relType, Props: props})
	return id, nil
}

// Neighbors returns nodes connected by the given relationship type.
func (g *Graph) Neighbors(_ context.Context, id, relType string) ([]Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var out []Node
	for _, e := range g.edges {
		if e.Type != relType {
			continue
		}
		if e.From == id {
			if n, ok := g.nodes[e.To]; ok {
				out = append(out, n)
			}
		}
	}
	return out, nil
}
