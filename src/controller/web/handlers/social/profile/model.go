package profile

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
	GetSubscriptions(c *fiber.Ctx) error
	GetSubscribers(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}
