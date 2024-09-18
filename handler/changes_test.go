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

func TestHandleChange(t *testing.T) {
	t.Parallel()
	scenarios := []hookScenario{
		{
			name:   "invalid field",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object":"instagram", 
				"entry":[{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "not-supported",
							"value": {
								"ignored": true
							}
					}]
				}]
			}`),
			expectErr: gometawebhooks.ErrChangesFieldNotImplemented,
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
				}
			},
		},
		{
			name:   "handles many",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object":"instagram", 
				"entry":[{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "mentions",
							"value": {
								"media_id": "999"
							}
					},{ 
						"field": "story_insights",
						"value": {
							"media_id": "999",
							"exits": 1,
							"replies": 2,
							"reach": 3,
							"taps_forward": 4,
							"taps_back": 5,
							"impressions": 6
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
							MediaID: "999",
						},
					}, {
						Field: "story_insights",
						Value: handler.StoryInsights{
							MediaID:     "999",
							Exits:       1,
							Replies:     2,
							Reach:       3,
							TapsForward: 4,
							TapsBack:    5,
							Impressions: 6,
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramChangesHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("change")
						return nil
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"change": 2,
			},
		},
		{
			name:   "caption mention",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object":"instagram", 
				"entry":[{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "mentions",
							"value": {
								"media_id": "999"
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
							MediaID: "999",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramMentionHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("mention")
						return nil
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"mention": 1,
			},
		},
		{
			name:   "comment mention",
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
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramMentionHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("mention")
						return nil
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"mention": 1,
			},
		},
		{
			name:   "story insights",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object":"instagram", 
				"entry":[{
					"id":"123",
					"time":1569262486134,
					"changes":[{ 
							"field": "story_insights",
							"value": {
								"media_id": "999",
								"exits": 1,
								"replies": 2,
								"reach": 3,
								"taps_forward": 4,
								"taps_back": 5,
								"impressions": 6
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
						Field: "story_insights",
						Value: handler.StoryInsights{
							MediaID:     "999",
							Exits:       1,
							Replies:     2,
							Reach:       3,
							TapsForward: 4,
							TapsBack:    5,
							Impressions: 6,
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramStoryInsightsHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("storyInsights")
						return nil
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"storyInsights": 1,
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
							"field": "story_insights",
							"value": {
								"media_id": "999",
								"exits": 1,
								"replies": 2,
								"reach": 3,
								"taps_forward": 4,
								"taps_back": 5,
								"impressions": 6
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
						Field: "story_insights",
						Value: handler.StoryInsights{
							MediaID:     "999",
							Exits:       1,
							Replies:     2,
							Reach:       3,
							TapsForward: 4,
							TapsBack:    5,
							Impressions: 6,
						},
					}},
				}},
			},
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
