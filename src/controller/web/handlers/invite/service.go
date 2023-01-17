package invite

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/errors"
	"core-server/src/internal/invites"
	"core-server/src/pkg/logger"
	smallcache "core-server/src/pkg/small-cache"
	"database/sql"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewHandler(l *logger.Logger, s *storage.Storage) Servicer {
	return &Handler{
		Logger:  l,
		Storage: s,
	}
}

func (h *Handler) Create(c *fiber.Ctx) error {
	userLogin := smallcache.GetUserLogin(c)

	u, err := h.Storage.Users.FetchLogin(userLogin)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	// check user permissions
	if u.Role != "admin" {
		return errors.SendError(c, http.StatusForbidden, errors.ErrNotEnoughPermissions)
	}

	// create invite
	i := invites.NewInvite(u.Id)

	if err := h.Storage.Invites.Create(i); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	// send invite
	return c.JSON(fiber.Map{
		"invite": i.Invite,
	})
}

func (h *Handler) Use(c *fiber.Ctx) error {
	// get invite
	invite := c.Params("invite", "")

	// validate invite
	if invite == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrNotFoundInvite)
	}

	i, err := h.Storage.Invites.Fetch(invite)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusBadRequest, errors.ErrNotFoundInvite)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if i.ActivatedBy != 0 {
		return errors.SendError(c, http.StatusForbidden, errors.ErrInviteIsAlreadyActivated)
	}

	// set invite to http only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "invite",
		Value:    i.Invite,
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(time.Hour * 720),
		HTTPOnly: true,
		Secure:   true,
	})

	// send ok
	return c.SendStatus(http.StatusOK)
}
