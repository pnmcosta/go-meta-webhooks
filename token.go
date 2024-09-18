package gometawebhooks

import "errors"

var (
	ErrVerifyTokenFailed = errors.New("invalid verify_token value")
)

func (hooks Webhooks) VerifyToken(queryValues map[string]string) (string, error) {
	mode := queryValues["hub.mode"]
	token := queryValues["hub.verify_token"]
	challenge := queryValues["hub.challenge"]
	if mode != "subscribe" || token != hooks.token || challenge == "" {
		return "", ErrVerifyTokenFailed
	}
	return challenge, nil
}
