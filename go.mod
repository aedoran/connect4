module mem0-go

go 1.24.3

require (
	github.com/gofiber/fiber/v2 v2.52.0
	github.com/jackc/pgx/v5 v5.5.0
)

replace github.com/gofiber/fiber/v2 => ./internal/fiber

replace github.com/jackc/pgx/v5 => ./internal/pgx
