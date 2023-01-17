package users

import (
	"core-server/src/pkg/crypto"
	"errors"
	"strings"
	"time"
)

func NewUser(login string, password string) *User {
	salt := crypto.GenerateString(12)
	hashed, _ := crypto.HashString(password, salt)

	u := User{
		Id:           0,
		Login:        strings.ToLower(login),
		Nickname:     login,
		Role:         "",
		Password:     hashed,
		Salt:         salt,
		Balance:      0,
		RegisterDate: time.Now().UnixMilli(),
	}

	return &u
}

func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if hashed, _ := crypto.HashString(oldPassword, u.Salt); hashed != u.Password {
		return errors.New("old password is not correct")
	}

	u.Salt = crypto.GenerateString(12)
	u.Password, _ = crypto.HashString(newPassword, u.Salt)

	return nil
}
