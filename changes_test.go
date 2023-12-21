package gometawebhooks_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

func TestHandleChanges(t *testing.T) {
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
