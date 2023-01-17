package turnstile

import (
	"net/http"
	"net/url"

	"github.com/goccy/go-json"
)

func NewService(secret string) Servicer {
	return &Service{
		Secret: secret,
	}
}

func (s *Service) Verify(response string) (bool, error) {
	data := url.Values{
		"secret":   {s.Secret},
		"response": {response},
	}

	r, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", data)
	if err != nil {
		return false, err
	}

	res := Response{}

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return false, err
	}

	return res.Success, nil
}
