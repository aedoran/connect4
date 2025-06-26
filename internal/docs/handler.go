package docs

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
)

//go:embed spec/openapi.yaml
var spec []byte

// Register sets up the /docs route serving the OpenAPI spec.
func Register(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Type("yaml")
		return c.Send(spec)
	})
}
