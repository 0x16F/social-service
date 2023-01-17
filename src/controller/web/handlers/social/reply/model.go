package reply

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
	GetAll(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type CreateReplyRequest struct {
	MessageId int64  `json:"message_id"`
	Content   string `json:"content"`
}
