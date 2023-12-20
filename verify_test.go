package gometawebhooks_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

func TestVerify(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		method    string
		expected  string
		expectErr error
	}{
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
			name:     "verifies",
			url:      "/webhooks/meta/?hub.mode=subscribe&hub.verify_token=meta_app_webhook_token&hub.challenge=challenge_response",
			method:   http.MethodGet,
			expected: "challenge_response",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hook, err := gometawebhooks.NewWebhooks(
				gometawebhooks.Options.Token("meta_app_webhook_token"),
			)
			if err != nil {
				t.Fatal(err)
				return
			}

			result, err := hook.Verify(httptest.NewRequest(test.method, test.url, nil))

			if test.expectErr != nil {
				if err == nil {
					t.Errorf("Expected an error, but got none.")
				}

				if !errors.Is(err, test.expectErr) {
					t.Errorf("Expected error %v, but got %v.", test.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}
