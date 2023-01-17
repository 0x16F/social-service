package message

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
	Like(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type CreateMessageRequest struct {
	Wall    int    `json:"wall"`
	Content string `json:"content"`
}

type CreateMessageResponse struct {
	Id int64 `json:"id"`
}
