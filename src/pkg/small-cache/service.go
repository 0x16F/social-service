package smallcache

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	userLoginKey = "X-UserLogin"
	userIdKey    = "X-UserId"
)

func SetUserId(c *fiber.Ctx, id int) {
	c.Response().Header.Set(userIdKey, strconv.Itoa(id))
}

func GetUserId(c *fiber.Ctx) (int, error) {
	return strconv.Atoi(string(c.Response().Header.Peek(userIdKey)))
}

func SetUserLogin(c *fiber.Ctx, login string) {
	c.Response().Header.Set(userLoginKey, login)
}

func GetUserLogin(c *fiber.Ctx) string {
	return string(c.Response().Header.Peek(userLoginKey))
}
