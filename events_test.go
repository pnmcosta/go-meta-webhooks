package gometawebhooks_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

func TestHandle(t *testing.T) {
	scenarios := []hookScenario{
		{
			name:      "invalid method",
			method:    http.MethodGet,
			expectErr: gometawebhooks.ErrInvalidHTTPMethod,
		},
		{
			name:      "nil body",
			method:    http.MethodPost,
			expectErr: gometawebhooks.ErrReadBodyPayload,
		},
		{
			name:      "empty body",
			method:    http.MethodPost,
			body:      strings.NewReader(``),
			expectErr: gometawebhooks.ErrReadBodyPayload,
		},
		{
			name:      "malformed body",
			method:    http.MethodPost,
			body:      strings.NewReader(`{"object`),
			expectErr: gometawebhooks.ErrParsingPayload,
		},
		{
			name:   "missing signature",
			method: http.MethodPost,
			options: []gometawebhooks.Option{
				gometawebhooks.Options.Secret("very_secret"),
			},
			body:      strings.NewReader(`{}`),
			expectErr: gometawebhooks.ErrMissingHubSignatureHeader,
		},
		{
			name:   "invalid signature",
			method: http.MethodPost,
			headers: map[string]string{
				"x_hub_signature_256": "1",
			},
			options: []gometawebhooks.Option{
				gometawebhooks.Options.Secret("very_secret"),
			},
			body:      strings.NewReader(`{}`),
			expectErr: gometawebhooks.ErrHMACVerificationFailed,
		},
		{
			name:   "verifies signature noop",
			method: http.MethodPost,
			headers: map[string]string{
				"x_hub_signature_256": genHmac("very_secret", `{"object":"instagram", "entry":[]}`),
			},
			options: []gometawebhooks.Option{
				gometawebhooks.Options.Secret("very_secret"),
			},
			body: strings.NewReader(`{"object":"instagram", "entry":[]}`),
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry:  []gometawebhooks.Entry{},
			},
		},
		{
			name:      "invalid payload",
			method:    http.MethodPost,
			body:      strings.NewReader(`{}`),
			expectErr: gometawebhooks.ErrInvalidPayload,
		},
		{
			name:      "unsupported object",
			method:    http.MethodPost,
			body:      strings.NewReader(`{"object":"none", "entry":[]}`),
			expectErr: gometawebhooks.ErrParsingEvent,
		},
		{
			name:   "no entries noop",
			method: http.MethodPost,
			body:   strings.NewReader(`{"object":"instagram", "entry":[]}`),
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry:  []gometawebhooks.Entry{},
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.test(t, func(t *testing.T) {
			hooks, req := scenario.init(t)

			ctx := context.Background()

			result, err := hooks.Handle(ctx, req)

			scenario.assert(t, result, err)
		})
	}
}
