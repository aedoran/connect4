module mem0-go

go 1.24.3

require (
	github.com/gofiber/fiber/v2 v2.52.0
	github.com/jackc/pgx/v5 v5.5.0
	github.com/jrallison/go-workers v0.0.0
	mem0-go/internal/observability v0.0.0
)

replace github.com/gofiber/fiber/v2 => ./internal/fiber

replace github.com/jackc/pgx/v5 => ./internal/pgx

replace github.com/jrallison/go-workers => ./internal/workers

replace mem0-go/internal/observability => ./internal/observability
