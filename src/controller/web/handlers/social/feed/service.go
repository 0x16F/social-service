package feed

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/errors"
	"core-server/src/pkg/logger"
	smallcache "core-server/src/pkg/small-cache"
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

func (h *Handler) Get(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

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

	feedType := c.Query("type")
	if feedType == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectFeedType)
	}

	messages, err := h.Storage.Socials.Feed.Get(feedType, userId, limitInt, offsetInt)
	if err != nil {
		if err.Error() == errors.ErrIncorrectFeedType {
			return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectFeedType)
		}

		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundFeed)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	count, err := h.Storage.Socials.Feed.GetCount(feedType, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundFeed)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(fiber.Map{
		"feed":  &messages,
		"count": &count,
	})
}
