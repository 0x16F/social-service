package web

import (
	"core-server/src/controller/web/handlers"
	"core-server/src/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	Router   *fiber.App
	Handlers *handlers.Handlers
	Logger   *logger.Logger
}
