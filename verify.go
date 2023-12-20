package gometawebhooks

import (
	"fmt"
	"net/http"
)

var (
	ErrVerificationFailed = fmt.Errorf("invalid verify_token value: %w", ErrWebhooks)
)

// Verify Meta Webhooks GET requests, when subscribing on App dashboard to objects and fields.
func (hook Webhooks) Verify(r *http.Request) (string, error) {
	if r.Method != http.MethodGet {
		return "", ErrInvalidHTTPMethod
	}

	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")
	if mode != "subscribe" || token != hook.token || challenge == "" {
		return "", ErrVerificationFailed
	}
	return challenge, nil
}
