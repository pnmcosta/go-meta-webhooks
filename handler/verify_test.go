package handler_test

import (
	"net/http"
	"testing"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

func TestVerify(t *testing.T) {
	scenarios := []hookScenario{
		{
			name:      "invalid method",
			method:    http.MethodPost,
			expectErr: gometawebhooks.ErrInvalidHTTPMethod,
		},
		{
			name:      "invalid mode",
			method:    http.MethodGet,
			expectErr: gometawebhooks.ErrVerifyTokenFailed,
		},
		{
			name:      "invalid verify_token",
			url:       "/webhooks/meta/?hub.mode=subscribe",
			method:    http.MethodGet,
			expectErr: gometawebhooks.ErrVerifyTokenFailed,
		},
		{
			name:      "missing challenge",
			url:       "/webhooks/meta/?hub.mode=subscribe&hub.verify_token=123",
			method:    http.MethodGet,
			expectErr: gometawebhooks.ErrVerifyTokenFailed,
		},
		{
			name: "verifies",
			url:  "/webhooks/meta/?hub.mode=subscribe&hub.verify_token=meta_app_webhook_token&hub.challenge=challenge_response",
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.Token("meta_app_webhook_token"),
				}
			},
			method:   http.MethodGet,
			expected: "challenge_response",
		},
	}

	for _, scenario := range scenarios {
		scenario.test(t, func(t *testing.T) {
			hooks, req := scenario.setup(t)

			result, err := hooks.HandleVerify(req)

			scenario.assert(t, result, nil, err)
		})
	}
}
