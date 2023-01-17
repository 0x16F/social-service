package turnstile

type Service struct {
	Secret string `json:"secret"`
}

type Response struct {
	Success     bool     `json:"success"`
	ErrorCodes  []string `json:"error-codes"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
}

type Servicer interface {
	Verify(response string) (bool, error)
}
