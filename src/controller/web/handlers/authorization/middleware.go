package authorization

import (
	"core-server/src/controller/web/errors"
	"core-server/src/pkg/jwt"
	smallcache "core-server/src/pkg/small-cache"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
)

func (h *Handler) IsAuthorized(c *fiber.Ctx) error {
	whitelist := []Whitelist{
		{
			Path:    "/auth",
			Methods: []string{"POST", "GET"},
		},
		{
			Path:    "/invite",
			Methods: []string{"GET"},
		},
		{
			Path:    "/avatar",
			Methods: []string{"GET"},
		},
		{
			Path:    "/background",
			Methods: []string{"GET"},
		},
	}

	for _, e := range whitelist {
		if strings.Contains(c.Path(), e.Path) && slices.Contains(e.Methods, c.Method()) {
			return c.Next()
		}
	}

	// get access token
	accessString := string(c.Request().Header.Peek("Authorization"))

	// validate access token
	access, err := h.JWT.ParseAccess(accessString)
	if err != nil {
		if err == jwt.ErrExpired {
			return errors.SendError(c, http.StatusUnauthorized, errors.ErrExpiredAccess)
		}

		return errors.SendError(c, http.StatusUnauthorized, errors.ErrIncorrectAccess)
	}

	smallcache.SetUserId(c, access.Id)
	smallcache.SetUserLogin(c, access.Login)

	return c.Next()
}
