package graphql

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"mem0-go/internal/memory"
)

// Request represents a minimal GraphQL request payload.
type Request struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// Register sets up GraphQL routes on the given app using the service.
func Register(app *fiber.App, svc *memory.Service) {
	app.Get("/graphql", func(c *fiber.Ctx) error {
		if c.Method() != http.MethodPost {
			return c.Type("html").SendString(playgroundHTML)
		}
		var req Request
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		}
		q := req.Query
		switch {
		case strings.Contains(q, "search") && strings.Contains(q, "query"):
			vecAny, _ := req.Variables["vector"].([]interface{})
			vec := make([]float32, len(vecAny))
			for i, v := range vecAny {
				if f, ok := v.(float64); ok {
					vec[i] = float32(f)
				}
			}
			limitF, _ := req.Variables["limit"].(float64)
			res, err := svc.Search(c.Context(), vec, int(limitF))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			return c.JSON(fiber.Map{"data": fiber.Map{"search": res}})
		case strings.Contains(q, "upsertMemory"):
			userF, _ := req.Variables["userID"].(float64)
			content, _ := req.Variables["content"].(string)
			vecAny, _ := req.Variables["vector"].([]interface{})
			vec := make([]float32, len(vecAny))
			for i, v := range vecAny {
				if f, ok := v.(float64); ok {
					vec[i] = float32(f)
				}
			}
			id, err := svc.StoreMemory(c.Context(), int64(userF), content, vec)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			return c.JSON(fiber.Map{"data": fiber.Map{"upsertMemory": fiber.Map{"id": id}}})
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unknown operation"})
		}
	})
}

const playgroundHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>GraphQL Playground</title>
  <link rel="stylesheet" href="https://unpkg.com/graphql-playground-react/build/static/css/index.css" />
  <link rel="shortcut icon" href="https://unpkg.com/graphql-playground-react/build/favicon.png" />
  <script src="https://unpkg.com/graphql-playground-react/build/static/js/middleware.js"></script>
</head>
<body>
  <div id="root" />
  <script>window.addEventListener('load', function () { GraphQLPlayground.init(document.getElementById('root'), { endpoint: '/graphql' }) })</script>
</body>
</html>`
