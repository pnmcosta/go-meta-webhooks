package gometawebhooks_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

func TestHandleChange(t *testing.T) {
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
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "not-supported",
						Value: nil,
					}},
				}},
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
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "mentions",
						Value: gometawebhooks.MentionsFieldValue{
							MediaID: "999",
						},
					}, {
						Field: "story_insights",
						Value: gometawebhooks.StoryInsightsFieldValue{
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
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.ChangesHandler(testHandler{func() {
						scenario.trigger("change")
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
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "mentions",
						Value: gometawebhooks.MentionsFieldValue{
							MediaID: "999",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.InstagramMentionHandler(testHandler{func() {
						scenario.trigger("mention")
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
					gometawebhooks.Options.InstagramMentionHandler(testHandler{func() {
						scenario.trigger("mention")
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
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "story_insights",
						Value: gometawebhooks.StoryInsightsFieldValue{
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
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.InstagramStoryInsightsHandler(testHandler{func() {
						scenario.trigger("storyInsights")
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
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Changes: []gometawebhooks.Change{{
						Field: "story_insights",
						Value: gometawebhooks.StoryInsightsFieldValue{
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
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.InstagramStoryInsightsHandler(testHandler{func() {
						time.Sleep(scenario.timeout + 5)
						scenario.trigger("storyInsights")
					}}),
				}
			},
			expectErr:        context.DeadlineExceeded,
			expectedHandlers: map[string]int{"storyInsights": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.test(t, func(t *testing.T) {
			hooks, req := scenario.setup(t)

			ctx, cancel := context.WithTimeout(context.Background(), scenario.timeout)
			defer cancel()

			result, payload, err := hooks.Handle(ctx, req)

			scenario.assert(t, result, payload, err)
		})
	}
}
