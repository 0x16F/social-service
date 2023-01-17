package jwt

import (
	"time"

	gjwt "github.com/golang-jwt/jwt"
)

func NewService(secrets *Service) Servicer {
	return &Service{
		AccessSecret:  secrets.AccessSecret,
		RefreshSecret: secrets.RefreshSecret,
	}
}

func (s *Service) GenerateAccess(id int, login string) (string, error) {
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS512, Token{
		StandardClaims: gjwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Minute * 5).Unix(),
		},
		Id:    id,
		Login: login,
	})

	str, err := token.SignedString([]byte(s.AccessSecret))

	return str, err
}

func (s *Service) GenerateRefresh(id int, login string) (string, error) {
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS512, Token{
		StandardClaims: gjwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 30).Unix(),
		},
		Id:    id,
		Login: login,
	})

	str, err := token.SignedString([]byte(s.RefreshSecret))

	return str, err
}

func (s *Service) ParseAccess(token string) (*Token, error) {
	t, err := gjwt.ParseWithClaims(token, &Token{}, func(t *gjwt.Token) (interface{}, error) {
		return []byte(s.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Token); ok && t.Valid {
		if time.Since(time.Now()).Milliseconds() > claims.ExpiresAt {
			return nil, ErrExpired
		}

		return claims, nil
	}

	return nil, err
}

func (s *Service) ParseRefresh(token string) (*Token, error) {
	t, err := gjwt.ParseWithClaims(token, &Token{}, func(t *gjwt.Token) (interface{}, error) {
		return []byte(s.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Token); ok && t.Valid {
		if time.Since(time.Now()).Milliseconds() > claims.ExpiresAt {
			return nil, ErrExpired
		}

		return claims, nil
	}

	return nil, err
}

func (s *Service) TokenToJson(t *Token) *TokenJson {
	token := TokenJson{
		Id:    t.Id,
		Login: t.Login,
	}

	return &token
}
