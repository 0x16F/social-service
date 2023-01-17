package errors

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func SendError(c *fiber.Ctx, code int, message string) error {
	e := ErrorResponse{
		Error: message,
	}

	data, _ := json.Marshal(&e)
	return fiber.NewError(code, string(data))
}
