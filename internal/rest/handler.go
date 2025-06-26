package rest

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"mem0-go/internal/memory"
)

// createMemoryRequest represents the payload for creating a memory.
type createMemoryRequest struct {
	UserID  int64     `json:"userID"`
	Content string    `json:"content"`
	Vector  []float32 `json:"vector"`
}

// searchRequest represents the payload for searching memories.
type searchRequest struct {
	Vector []float32 `json:"vector"`
	Limit  int       `json:"limit"`
}

// Register sets up REST routes on the given app using the service.
func Register(app *fiber.App, svc *memory.Service) {
	// @Summary Create memory
	// @Description Store memory text and embedding
	// @Tags memories
	// @Accept json
	// @Produce json
	// @Param data body createMemoryRequest true "memory info"
	// @Success 200 {object} map[string]int64
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /api/v1/memories [post]
	app.Post("/api/v1/memories", func(c *fiber.Ctx) error {
		var req createMemoryRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		}
		id, err := svc.StoreMemory(c.Context(), req.UserID, req.Content, req.Vector)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"id": id})
	})

	// @Summary Search memories
	// @Description Semantic search over stored memories
	// @Tags memories
	// @Accept json
	// @Produce json
	// @Param data body searchRequest true "search parameters"
	// @Success 200 {object} map[string][]memory.MemoryResult
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /api/v1/memories/search [post]
	app.Post("/api/v1/memories/search", func(c *fiber.Ctx) error {
		var req searchRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		}
		res, err := svc.Search(c.Context(), req.Vector, req.Limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"results": res})
	})

	// @Summary Get memory
	// @Description Retrieve memory by ID
	// @Tags memories
	// @Produce json
	// @Param id path int true "Memory ID"
	// @Success 200 {object} db.Memory
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /api/v1/memories/{id} [get]
	app.Get("/api/v1/memories/", func(c *fiber.Ctx) error {
		parts := strings.Split(c.Path(), "/")
		if len(parts) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		idStr := parts[len(parts)-1]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		m, err := svc.GetMemory(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(m)
	})
}
