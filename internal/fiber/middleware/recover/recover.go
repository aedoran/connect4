package recover

import "github.com/gofiber/fiber/v2"

// New returns a no-op recovery middleware compatible with the fiber interface.
func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
