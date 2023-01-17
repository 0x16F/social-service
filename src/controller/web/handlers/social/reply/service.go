package reply

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/errors"
	"core-server/src/internal/social/reply"
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

// POST social/replies
func (h *Handler) Create(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

	u, err := h.Storage.Users.Fetch(userId)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	// parse request body
	request := CreateReplyRequest{}

	if err := c.BodyParser(&request); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectBody)
	}

	// is message exists
	message, err := h.Storage.Socials.Messages.Get(request.MessageId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundMessage)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	// validate content
	if len(strings.Join(strings.Fields(request.Content), "")) == 0 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectContent)
	}

	// is profile exists
	profile, err := h.Storage.Socials.Profiles.Get(userId, message.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundUser)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if profile.IsCommentsClosed && u.Id != profile.UserId {
		return errors.SendError(c, http.StatusForbidden, errors.ErrClosedComments)
	}

	reply := reply.Reply{
		Id:        time.Now().UnixMilli(),
		ProfileId: profile.UserId,
		MessageId: request.MessageId,
		Content:   request.Content,
		Author:    u.Id,
		Nickname:  u.Nickname,
	}

	// create reply
	if err := h.Storage.Socials.Replies.Create(&reply); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(&reply)
}

// DELETE social/replies/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userLogin := smallcache.GetUserLogin(c)

	u, err := h.Storage.Users.FetchLogin(userLogin)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	// get reply id from params
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	// get reply
	reply, err := h.Storage.Socials.Replies.Get(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundReply)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	// check permissions
	if u.Id != reply.Author && u.Role != "admin" && u.Id != reply.ProfileId {
		return errors.SendError(c, http.StatusForbidden, errors.ErrNotEnoughPermissions)
	}

	// delete reply
	if err := h.Storage.Socials.Replies.Delete(id); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.SendStatus(http.StatusOK)
}

// GET social/replies/:id
func (h *Handler) GetAll(c *fiber.Ctx) error {
	// get reply id from params
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	if _, err := h.Storage.Socials.Messages.Get(id); err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundMessage)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	replies, err := h.Storage.Socials.Replies.GetAll(id)
	if err != nil {
		if err != sql.ErrNoRows {
			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}
	}

	return c.JSON(&fiber.Map{
		"replies": &replies,
	})
}
