package gometawebhooks

import "fmt"

var (
	ErrVerifyTokenFailed = fmt.Errorf("invalid verify_token value: %w", ErrWebhooks)
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
