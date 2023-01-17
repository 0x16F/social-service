package jwt

import (
	"errors"

	gjwt "github.com/golang-jwt/jwt"
)

type Token struct {
	gjwt.StandardClaims
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type TokenJson struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type Service struct {
	AccessSecret  string `json:"access"`
	RefreshSecret string `json:"refresh"`
}

type Servicer interface {
	ParseAccess(token string) (*Token, error)
	ParseRefresh(token string) (*Token, error)
	GenerateAccess(id int, login string) (string, error)
	GenerateRefresh(id int, login string) (string, error)
	TokenToJson(token *Token) *TokenJson
}

var (
	ErrExpired = errors.New("expired")
)
