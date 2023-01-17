package message

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/errors"
	"core-server/src/internal/social/like"
	"core-server/src/internal/social/message"
	"core-server/src/pkg/logger"
	smallcache "core-server/src/pkg/small-cache"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewHandler(l *logger.Logger, s *storage.Storage) Servicer {
	return &Handler{
		Logger:  l,
		Storage: s,
	}
}

// GET social/messages/:id?limit=&offset=
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

	// get user id from params
	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

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

	// fetch messages from wall
	messages, err := h.Storage.Socials.Messages.GetAll(id, userId, offsetInt, limitInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundMessage)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	count, err := h.Storage.Socials.Messages.GetCount(id)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(&fiber.Map{
		"messages": &messages,
		"count":    &count,
	})
}

// POST social/messages
func (h *Handler) Create(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)
	userLogin := smallcache.GetUserLogin(c)

	request := CreateMessageRequest{}

	if err := c.BodyParser(&request); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectBody)
	}

	if len(strings.Join(strings.Fields(request.Content), "")) == 0 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectContent)
	}

	u, err := h.Storage.Users.Fetch(request.Wall)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundUser)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	wall, err := h.Storage.Socials.Profiles.Get(userId, u.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundUser)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if wall.IsWallClosed && !strings.EqualFold(userLogin, u.Login) {
		return errors.SendError(c, http.StatusForbidden, errors.ErrClosedWall)
	}

	message := message.Message{
		Id:       time.Now().UnixMilli(),
		AuthorId: userId,
		Content:  request.Content,
		UserId:   request.Wall,
	}

	if err := h.Storage.Socials.Messages.Create(&message); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(&CreateMessageResponse{
		Id: message.Id,
	})
}

// DELETE social/messages/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userLogin := smallcache.GetUserLogin(c)

	u, err := h.Storage.Users.FetchLogin(userLogin)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	// get message id from params
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	message, err := h.Storage.Socials.Messages.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundMessage)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if message.UserId != u.Id && message.AuthorId != u.Id && u.Role != "admin" {
		return errors.SendError(c, http.StatusForbidden, errors.ErrNotEnoughPermissions)
	}

	if err := h.Storage.Socials.Messages.Delete(id); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.SendStatus(http.StatusOK)
}

// POST social/messages/:id
func (h *Handler) Like(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

	// get message id from params
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	// is message exists
	message, err := h.Storage.Socials.Messages.GetWithInfo(int64(id), userId)
	if err != nil {
		h.Logger.Error(err)

		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundMessage)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	// is message already liked
	if message.IsLiked {
		if err := h.Storage.Socials.Likes.Delete(&like.Like{
			Id:   int64(id),
			From: userId,
		}); err != nil {
			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}
	} else {
		if err := h.Storage.Socials.Likes.Create(&like.Like{
			Id:   int64(id),
			From: userId,
		}); err != nil {
			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}
	}

	// fetch message with new params
	message, err = h.Storage.Socials.Messages.GetWithInfo(int64(id), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundMessage)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(&message)
}
