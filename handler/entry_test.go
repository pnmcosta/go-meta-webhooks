package handler_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
	"github.com/pnmcosta/go-meta-webhooks/handler"
)

func TestHandleEvent(t *testing.T) {
	t.Parallel()
	scenarios := []hookScenario{
		{
			name:      "invalid method",
			method:    http.MethodGet,
			expectErr: handler.ErrInvalidHTTPMethod,
		},
		{
			name:      "nil body",
			method:    http.MethodPost,
			expectErr: handler.ErrReadBodyPayload,
		},
		{
			name:      "empty body",
			method:    http.MethodPost,
			body:      strings.NewReader(``),
			expectErr: handler.ErrReadBodyPayload,
		},
		{
			name:      "malformed body",
			method:    http.MethodPost,
			body:      strings.NewReader(`{"object`),
			expectErr: gometawebhooks.ErrParsingPayload,
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
				}
			},
		},
		{
			name:   "missing signature",
			method: http.MethodPost,
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.Secret("very_secret"),
				}
			},
			body:      strings.NewReader(`{}`),
			expectErr: gometawebhooks.ErrMissingHubSignatureHeader,
		},
		{
			name:   "invalid signature",
			method: http.MethodPost,
			headers: map[string]string{
				"X-Hub-Signature-256": "1",
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.Secret("very_secret"),
				}
			},
			body:      strings.NewReader(`{}`),
			expectErr: gometawebhooks.ErrHMACVerificationFailed,
		},
		{
			name:   "verifies signature noop",
			method: http.MethodPost,
			headers: map[string]string{
				"X-Hub-Signature-256": genHmac("very_secret", `{"object":"instagram", "entry":[]}`),
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.Secret("very_secret"),
				}
			},
			body: strings.NewReader(`{"object":"instagram", "entry":[]}`),
			expected: handler.Event{
				Object: handler.Instagram,
				Entry:  []handler.Entry{},
			},
		},
		{
			name:      "invalid payload",
			method:    http.MethodPost,
			body:      strings.NewReader(`{}`),
			expectErr: gometawebhooks.ErrInvalidPayload,
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
				}
			},
		},
		{
			name:   "handles unsupported object entries",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object":"unsupported", 
				"entry":[{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "mentions",
							"value": {
								"media_id": "999",
								"comment_id": "4444"
							}
					}]
				}]
			}`),
			expected: handler.Event{
				Object: "unsupported",
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []handler.Change{{
						Field: "mentions",
						Value: handler.Mention{
							MediaID:   "999",
							CommentID: "4444",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
				}
			},
			expectErr: gometawebhooks.ErrObjectNotSupported,
		},
		{
			name:   "no entries noop",
			method: http.MethodPost,
			body:   strings.NewReader(`{"object":"instagram", "entry":[]}`),
			expected: handler.Event{
				Object: handler.Instagram,
				Entry:  []handler.Entry{},
			},
		},
		{
			name:   "handles many entries",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object":"instagram", 
				"entry":[{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "mentions",
							"value": {
								"media_id": "999",
								"comment_id": "4444"
							}
					}]
				},{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "mentions",
							"value": {
								"media_id": "999",
								"comment_id": "4444"
							}
					}]
				},{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "mentions",
							"value": {
								"media_id": "999",
								"comment_id": "4444"
							}
					}]
				}]
				}`),
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []handler.Change{{
						Field: "mentions",
						Value: handler.Mention{
							MediaID:   "999",
							CommentID: "4444",
						},
					}},
				}, {
					Id:   "123",
					Time: 1569262486134,
					Changes: []handler.Change{{
						Field: "mentions",
						Value: handler.Mention{
							MediaID:   "999",
							CommentID: "4444",
						},
					}},
				}, {
					Id:   "123",
					Time: 1569262486134,
					Changes: []handler.Change{{
						Field: "mentions",
						Value: handler.Mention{
							MediaID:   "999",
							CommentID: "4444",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("entry")
						return nil
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"entry": 3,
			},
		},
		{
			name:   "deadline exceeded",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object":"instagram", 
				"entry":[{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "mentions",
							"value": {
								"media_id": "999",
								"comment_id": "4444"
							}
					}]
				}]
			}`),
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramHandler(testHandler{func(ctx context.Context) error {
						time.Sleep(scenario.timeout * 2)
						return context.Cause(ctx)
					}}),
				}
			},
			expectErr: context.DeadlineExceeded,
			timeout:   50 * time.Millisecond,
		},
	}

	for _, scenario := range scenarios {
		scenario.test(t, func(t *testing.T) {
			hooks, req := scenario.setup(t)

			ctx, cancel := context.WithTimeout(context.Background(), scenario.timeout)
			defer cancel()

			result, payload, err := hooks.HandleRequest(ctx, req)

			scenario.assert(t, result, payload, err)
		})
	}
}
