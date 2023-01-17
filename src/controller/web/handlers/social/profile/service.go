package profile

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/errors"
	"core-server/src/pkg/logger"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func NewHandler(l *logger.Logger, s *storage.Storage) Servicer {
	return &Handler{
		Logger:  l,
		Storage: s,
	}
}

// GET social/profiles?nickname=&limit=&offset=
func (h *Handler) GetAll(c *fiber.Ctx) error {
	nickname := c.Query("nickname")

	limit := c.Query("limit")
	if limit == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLimit)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLimit)
	}

	if limitInt > 30 || limitInt <= 0 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrLimitValueErr)
	}

	offset := c.Query("offset")

	if offset == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectOffset)
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectOffset)
	}

	if nickname != "" {
		// fetch messages from wall
		profiles, err := h.Storage.Socials.Profiles.GetAllLike(nickname, offsetInt, limitInt)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundMessage)
			}

			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}

		return c.JSON(&fiber.Map{
			"profiles": &profiles,
		})
	}

	// fetch messages from wall
	profiles, err := h.Storage.Socials.Profiles.GetAll(offsetInt, limitInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundMessage)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	count, err := h.Storage.Socials.Profiles.GetCount()
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(&fiber.Map{
		"profiles": &profiles,
		"count":    &count,
	})
}

// GET social/:id/subscribers?nickname=&limit=&offset=
func (h *Handler) GetSubscribers(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	nickname := c.Query("nickname")

	limit := c.Query("limit")
	if limit == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLimit)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLimit)
	}

	if limitInt > 30 || limitInt <= 0 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrLimitValueErr)
	}

	offset := c.Query("offset")
	if offset == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectOffset)
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectOffset)
	}

	if nickname != "" {
		subscribers, err := h.Storage.Socials.Subs.GetSubscribersLike(nickname, id, offsetInt, limitInt)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundSubscribers)
			}

			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}

		return c.JSON(fiber.Map{
			"subscribers": &subscribers,
		})
	}

	subscribers, err := h.Storage.Socials.Subs.GetSubscribers(id, offsetInt, limitInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundSubscribers)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(fiber.Map{
		"subscribers": &subscribers,
	})
}

// GET social/:id/subscriptions?nickname=&limit=&offset=
func (h *Handler) GetSubscriptions(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	nickname := c.Query("nickname")

	limit := c.Query("limit")
	if limit == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLimit)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLimit)
	}

	if limitInt > 30 || limitInt <= 0 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrLimitValueErr)
	}

	offset := c.Query("offset")
	if offset == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectOffset)
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectOffset)
	}

	if nickname != "" {
		subscriptions, err := h.Storage.Socials.Subs.GetSubscriptionsLike(nickname, id, offsetInt, limitInt)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundSubscriptions)
			}

			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}

		return c.JSON(fiber.Map{
			"subscriptions": &subscriptions,
		})
	}

	subscriptions, err := h.Storage.Socials.Subs.GetSubscriptions(id, offsetInt, limitInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundSubscriptions)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(fiber.Map{
		"subscriptions": &subscriptions,
	})
}
