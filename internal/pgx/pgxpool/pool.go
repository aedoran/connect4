package pgxpool

import "context"

type Pool struct{}

type Row struct{}

type CommandTag struct{}

type Rows struct{}

func New(ctx context.Context, connString string) (*Pool, error) {
	return &Pool{}, nil
}

func (p *Pool) Close() {}

func (p *Pool) Ping(ctx context.Context) error { return nil }

func (p *Pool) Exec(ctx context.Context, sql string, args ...interface{}) (CommandTag, error) {
	return CommandTag{}, nil
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	return Row{}
}

func (r Row) Scan(dest ...interface{}) error { return nil }

func (p *Pool) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	return Rows{}, nil
}

func (r Rows) Close() {}

func (r Rows) Next() bool { return false }

func (r Rows) Scan(dest ...interface{}) error { return nil }
