package graph

import (
	"context"
	"net"
	"testing"
)

func TestCreateAndQuery(t *testing.T) {
	// stub dialer to avoid real network
        dialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
                c1, c2 := net.Pipe()
                go func() {
                        _ = c1.Close()
                }()
                return c2, nil
        }

	cfg := LoadConfig()
	g, err := Connect(context.Background(), cfg)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}

	n1, _ := g.CreateNode(context.Background(), "Person", nil)
	n2, _ := g.CreateNode(context.Background(), "Person", nil)
	if _, err := g.CreateEdge(context.Background(), n1, n2, "KNOWS", nil); err != nil {
		t.Fatalf("create edge: %v", err)
	}

	neigh, err := g.Neighbors(context.Background(), n1, "KNOWS")
	if err != nil {
		t.Fatalf("neighbors: %v", err)
	}
	if len(neigh) != 1 || neigh[0].ID != n2 {
		t.Fatalf("unexpected neighbors: %+v", neigh)
	}
}
