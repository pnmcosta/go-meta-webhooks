package gometawebhooks_test

import (
	"net/http"
	"testing"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

func TestVerify(t *testing.T) {
	scenarios := []hookScenario{
		{
			name:      "invalid method",
			url:       "/webhooks/meta",
			method:    http.MethodPost,
			expectErr: gometawebhooks.ErrInvalidHTTPMethod,
		},
		{
			name:      "invalid mode",
			url:       "/webhooks/meta",
			method:    http.MethodGet,
			expectErr: gometawebhooks.ErrVerificationFailed,
		},
		{
			name:      "invalid verify_token",
			url:       "/webhooks/meta/?hub.mode=subscribe",
			method:    http.MethodGet,
			expectErr: gometawebhooks.ErrVerificationFailed,
		},
		{
			name:      "missing challenge",
			url:       "/webhooks/meta/?hub.mode=subscribe&hub.verify_token=123",
			method:    http.MethodGet,
			expectErr: gometawebhooks.ErrVerificationFailed,
		},
		{
			name: "verifies",
			url:  "/webhooks/meta/?hub.mode=subscribe&hub.verify_token=meta_app_webhook_token&hub.challenge=challenge_response",
			options: []gometawebhooks.Option{
				gometawebhooks.Options.Token("meta_app_webhook_token"),
			},
			method:   http.MethodGet,
			expected: "challenge_response",
		},
	}

	for _, scenario := range scenarios {
		scenario.test(t, func(t *testing.T) {
			hooks, req := scenario.init(t)

			result, err := hooks.Verify(req)

			scenario.assert(t, result, err)
		})
	}
}
