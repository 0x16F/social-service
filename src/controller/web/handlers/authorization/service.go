package authorization

import (
	"core-server/src/controller/storage"
	"core-server/src/controller/web/errors"
	"core-server/src/internal/invites"
	"core-server/src/internal/users"
	"core-server/src/pkg/crypto"
	"core-server/src/pkg/jwt"
	"core-server/src/pkg/logger"
	"core-server/src/pkg/turnstile"
	"database/sql"
	"net/http"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewHandler(l *logger.Logger, j jwt.Servicer, s *storage.Storage, c turnstile.Servicer) Servicer {
	return &Handler{
		Logger:  l,
		Storage: s,
		JWT:     j,
		CF:      c,
	}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	request := LoginRequest{}

	if err := c.BodyParser(&request); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectBody)
	}

	h.Logger.Infof("LOGIN: %s", request.Login)

	verified, err := h.CF.Verify(request.Solution)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if !verified {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrBadCaptcha)
	}

	// is login valid
	if !regexp.MustCompile(`(?m)^[A-Za-z0-9]+$`).Match([]byte(request.Login)) {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLogin)
	}

	if len(request.Login) < 3 || len(request.Login) > 24 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLoginLength)
	}

	// is password valid
	passwordLength := len(request.Password)

	if passwordLength == 0 || passwordLength < 8 || passwordLength > 64 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectPasswordLength)
	}

	u, err := h.Storage.Users.FetchLogin(request.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundUser)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	hashed, err := crypto.HashString(request.Password, u.Salt)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if hashed != u.Password {
		return errors.SendError(c, http.StatusForbidden, errors.ErrLoginFailed)
	}

	refresh, err := h.JWT.GenerateRefresh(u.Id, u.Login)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	access, err := h.JWT.GenerateAccess(u.Id, u.Login)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    refresh,
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(time.Hour * 720),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{
		"token": access,
	})
}

func (h *Handler) Register(c *fiber.Ctx) error {
	request := RegisterRequest{}

	if err := c.BodyParser(&request); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectBody)
	}

	h.Logger.Infof("REGISTER: %s", request.Login)

	verified, err := h.CF.Verify(request.Solution)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if !verified {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrBadCaptcha)
	}

	invite := c.Cookies("invite", "")

	if invite == "" {
		return errors.SendError(c, http.StatusForbidden, errors.ErrRegisterInviteOnly)
	}

	fetchedInvite, err := h.Storage.Invites.Fetch(invite)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.SendError(c, http.StatusNotFound, errors.ErrNotFoundInvite)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	if fetchedInvite.ActivatedBy != 0 {
		return errors.SendError(c, http.StatusForbidden, errors.ErrInviteIsAlreadyActivated)
	}

	// is login valid
	if !regexp.MustCompile(`(?m)^[A-Za-z0-9]+$`).Match([]byte(request.Login)) {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLogin)
	}

	if len(request.Login) < 3 || len(request.Login) > 24 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectLoginLength)
	}

	// is password valid
	passwordLength := len(request.Password)

	if passwordLength == 0 || passwordLength < 8 || passwordLength > 64 {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrIncorrectPasswordLength)
	}

	if u, err := h.Storage.Users.FetchLogin(request.Login); err != nil || u != nil {
		if u != nil {
			return errors.SendError(c, http.StatusForbidden, errors.ErrAccountIsAlreadyExists)
		}

		if err != sql.ErrNoRows {
			h.Logger.Error(err)
			return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
		}
	}

	u := users.NewUser(request.Login, request.Password)

	id, err := h.Storage.Users.Create(u)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	activatedInvite := &invites.ActivatedInvite{
		Invite:       fetchedInvite.Invite,
		ActivatedBy:  *id,
		ActivateTime: time.Now().UnixMilli(),
	}

	if err := h.Storage.Invites.Activate(activatedInvite); err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	u, err = h.Storage.Users.FetchLogin(u.Login)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	refresh, err := h.JWT.GenerateRefresh(u.Id, u.Login)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	access, err := h.JWT.GenerateAccess(u.Id, u.Login)
	if err != nil {
		h.Logger.Error(err)
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    refresh,
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(time.Hour * 720),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{
		"token": access,
	})
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	// get refresh token
	refreshString := c.Cookies("refresh", "")

	if refreshString == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrNotFoundRefresh)
	}

	// is refresh token valid
	refresh, err := h.JWT.ParseRefresh(refreshString)
	if err != nil {
		c.Cookie(&fiber.Cookie{
			Name:     "refresh",
			Value:    "",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HTTPOnly: true,
		})

		if err == jwt.ErrExpired {
			return errors.SendError(c, http.StatusUnauthorized, errors.ErrExpiredRefresh)
		}

		h.Logger.Error(err)
		return errors.SendError(c, http.StatusForbidden, errors.ErrIncorrectRefresh)
	}

	// generate access
	access, err := h.JWT.GenerateAccess(refresh.Id, refresh.Login)
	if err != nil {
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	refreshString, err = h.JWT.GenerateRefresh(refresh.Id, refresh.Login)
	if err != nil {
		return errors.SendError(c, http.StatusInternalServerError, errors.ErrInternalError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    refreshString,
		MaxAge:   30 * 24 * 60 * 60,
		Expires:  time.Now().Add(time.Hour * 720),
		HTTPOnly: true,
		Secure:   true,
	})

	// send access
	return c.JSON(fiber.Map{
		"token": access,
	})
}

func (h *Handler) ParseToken(c *fiber.Ctx) error {
	// get access
	accessString := c.Params("token", "")

	if accessString == "" {
		return errors.SendError(c, http.StatusBadRequest, errors.ErrNotFoundAccess)
	}

	// is access token valid
	access, err := h.JWT.ParseAccess(accessString)
	if err != nil {
		if err == jwt.ErrExpired {
			return errors.SendError(c, http.StatusUnauthorized, errors.ErrExpiredAccess)
		}

		return errors.SendError(c, http.StatusForbidden, errors.ErrIncorrectAccess)
	}

	token := h.JWT.TokenToJson(access)

	return c.JSON(&token)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
	})

	return c.SendStatus(http.StatusOK)
}
