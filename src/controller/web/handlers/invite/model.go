package invite

import (
	"core-server/src/controller/storage"
	"core-server/src/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Logger  *logger.Logger
	Storage *storage.Storage
}

type Servicer interface {
	Create(c *fiber.Ctx) error
	Use(c *fiber.Ctx) error
}
