package handler_test

import (
	"net/http"
	"testing"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
	"github.com/pnmcosta/go-meta-webhooks/handler"
)

func TestVerify(t *testing.T) {
	t.Parallel()
	scenarios := []hookScenario{
		{
			name:      "invalid method",
			method:    http.MethodPost,
			expectErr: handler.ErrInvalidHTTPMethod,
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
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.Token("meta_app_webhook_token"),
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
