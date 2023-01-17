package user

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/errors"
	"core-server/src/pkg/logger"
	smallcache "core-server/src/pkg/small-cache"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

func NewHandler(l *logger.Logger, s *storage.Storage) *Handler {
	return &Handler{
		Logger:  l,
		Storage: s,
	}
}

// Change password
// Responses
// 200 - ok
// 400 - failed to parse json
// 401 - token is not valid or expired
// 403 - old password is not valid
// 404 - user is not found
// 500 - internal error
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	userLogin := smallcache.GetUserLogin(c)

	request := ChangePasswordRequest{}

	if err := c.BodyParser(&request); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectBody)
	}

	u, err := h.Storage.Users.FetchLogin(userLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundUser)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if err := u.ChangePassword(request.OldPassword, request.NewPassword); err != nil {
		return errors.SendError(c, http.StatusForbidden, err.Error())
	}

	if err := h.Storage.Users.Update(u); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.SendStatus(http.StatusOK)
}

// get my profile info
// Responses
// 200 - ok
// 401 - token is not valid or expired
// 404 - user is not found
// 500 - internal error
func (h *Handler) GetMyInfo(c *fiber.Ctx) error {
	userLogin := smallcache.GetUserLogin(c)

	u, err := h.Storage.Users.FetchLogin(userLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundUser)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	response := UserInfoResponse{
		Id:       u.Id,
		Nickname: u.Nickname,
		Balance:  u.Balance,
		Role:     u.Role,
	}

	return c.JSON(&response)
}

func (h *Handler) GetProfileInfo(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

	id, err := c.ParamsInt("id", -1)
	if err != nil || id == -1 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	profile, err := h.Storage.Socials.Profiles.Get(userId, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusBadRequest, errors.ErrNotFoundUser)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.JSON(&profile)
}

func (h *Handler) GetUserAvatar(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	if _, err := os.Stat(fmt.Sprintf("./avatars/%d/%d.png", id, id)); err != nil {
		return c.SendFile("./avatars/no_avatar.png")
	}

	return c.SendFile(fmt.Sprintf("./avatars/%d/%d.png", id, id))
}

func (h *Handler) SetUserAvatar(c *fiber.Ctx) error {
	userId, err := smallcache.GetUserId(c)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if _, err := os.Stat(fmt.Sprintf("./avatars/%d", userId)); err != nil {
		if err := os.Mkdir(fmt.Sprintf("./avatars/%d", userId), 0777); err != nil {
			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}
	}

	file, err := c.FormFile("image")
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrNotFoundImage)
	}

	if file.Size > 10<<20 {
		return errors.SendError(c, http.StatusForbidden, errors.ErrImageSizeIsTooBig)
	}

	c.SaveFile(file, fmt.Sprintf("./avatars/%d/%d.png", userId, userId))

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) DeleteUserAvatar(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

	path := fmt.Sprintf("./avatars/%d/%d.png", userId, userId)

	if _, err := os.Stat(path); err != nil {
		c.SendStatus(http.StatusOK)
	}

	if err := os.Remove(path); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) GetUserBackground(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	if _, err := os.Stat(fmt.Sprintf("./avatars/%d/%d_background.png", id, id)); err != nil {
		return c.SendFile("./avatars/default_background.png")
	}

	return c.SendFile(fmt.Sprintf("./avatars/%d/%d_background.png", id, id))
}

func (h *Handler) SetUserBackground(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

	if _, err := os.Stat(fmt.Sprintf("./avatars/%d", userId)); err != nil {
		if err := os.Mkdir(fmt.Sprintf("./avatars/%d", userId), 0777); err != nil {
			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}
	}

	file, err := c.FormFile("image")
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrNotFoundImage)
	}

	if file.Size > 10<<20 {
		return errors.SendError(c, http.StatusForbidden, errors.ErrImageSizeIsTooBig)
	}

	c.SaveFile(file, fmt.Sprintf("./avatars/%d/%d_background.png", userId, userId))

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) DeleteUserBackground(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

	path := fmt.Sprintf("./avatars/%d/%d_background.png", userId, userId)

	if _, err := os.Stat(path); err != nil {
		c.SendStatus(http.StatusOK)
	}

	if err := os.Remove(path); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	return c.SendStatus(http.StatusOK)
}

// POST user/:id/subscribe
func (h *Handler) Subscribe(c *fiber.Ctx) error {
	userId, _ := smallcache.GetUserId(c)

	id, err := c.ParamsInt("id")
	if err != nil {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectId)
	}

	if userId == id {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrCantSubscribeYourself)
	}

	u, err := h.Storage.Socials.Profiles.Get(userId, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundUser)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if !u.IsSubscribed {
		if err := h.Storage.Socials.Subs.Subscribe(userId, id); err != nil {
			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}

		u.Subscribers += 1
	}

	if u.IsSubscribed {
		if err := h.Storage.Socials.Subs.Unsubscribe(userId, id); err != nil {
			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}

		u.Subscribers -= 1
	}

	u.IsSubscribed = !u.IsSubscribed

	return c.JSON(&u)
}
