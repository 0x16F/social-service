package user

import (
	"core-server/src/controller/storage"
	"core-server/src/pkg/logger"
)

type Handler struct {
	Logger  *logger.Logger
	Storage *storage.Storage
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserInfoResponse struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Balance  int64  `json:"balance"`
	Role     string `json:"role"`
}

type SetAvatarRequest struct {
	Image string `form:"image"`
}
