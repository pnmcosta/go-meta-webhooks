package gometawebhooks_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

func TestHandleEvent(t *testing.T) {
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
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.Secret("very_secret"),
				}
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
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.Secret("very_secret"),
				}
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
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.Secret("very_secret"),
				}
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
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "mentions",
						Value: gometawebhooks.MentionsFieldValue{
							MediaID:   "999",
							CommentID: "4444",
						},
					}},
				}, {
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "mentions",
						Value: gometawebhooks.MentionsFieldValue{
							MediaID:   "999",
							CommentID: "4444",
						},
					}},
				}, {
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "mentions",
						Value: gometawebhooks.MentionsFieldValue{
							MediaID:   "999",
							CommentID: "4444",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.EntryHandler(testHandler{func() {
						scenario.trigger("entry")
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
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "mentions",
						Value: gometawebhooks.MentionsFieldValue{
							MediaID:   "999",
							CommentID: "4444",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.EntryHandler(testHandler{func() {
						time.Sleep(scenario.timeout + 1)
						scenario.trigger("entry")
					}}),
				}
			},
			expectErr:        context.DeadlineExceeded,
			expectedHandlers: map[string]int{"entry": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.test(t, func(t *testing.T) {
			hooks, req := scenario.setup(t)

			ctx, cancel := context.WithTimeout(context.Background(), scenario.timeout)
			defer cancel()

			result, err := hooks.Handle(ctx, req)

			scenario.assert(t, result, err)
		})
	}
}
