package gometawebhooks_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

func TestHandleMessaging(t *testing.T) {
	scenarios := []hookScenario{
		{
			name:   "handles many",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object": "instagram",
				"entry": [
				  {
					"id": "123",
					"time": 1569262486134,
					"messaging": [
					  {
						"sender": {
						  "id": "567"
						},
						"recipient": {
						  "id": "123"
						},
						"timestamp": 1569262485349,
						"message":{
							"mid": "MESSAGE_ID",
							"text": "message from 567"
						}
					  },
					  {
						"sender": {
						  "id": "444"
						},
						"recipient": {
						  "id": "123"
						},
						"timestamp": 1569262485349,
						"message":{
							"mid": "MESSAGE_ID",
							"text": "message from 444"
						}
					  },
					  {
						"sender": {
						  "id": "666"
						},
						"recipient": {
						  "id": "123"
						},
						"timestamp": 1569262485349,
						"message":{
							"mid": "MESSAGE_ID",
							"text": "message from 666"
						}
					  }
					]
				  }
				]
			  }`),
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []gometawebhooks.Messaging{{
						Sender: struct {
							Id string "json:\"id\""
						}{
							Id: "567",
						},
						Recipient: struct {
							Id string "json:\"id\""
						}{
							Id: "123",
						},
						Timestamp: 1569262485349,
						Message: gometawebhooks.Message{
							Id:   "MESSAGE_ID",
							Text: "message from 567",
						},
					}, {
						Sender: struct {
							Id string "json:\"id\""
						}{
							Id: "444",
						},
						Recipient: struct {
							Id string "json:\"id\""
						}{
							Id: "123",
						},
						Timestamp: 1569262485349,
						Message: gometawebhooks.Message{
							Id:   "MESSAGE_ID",
							Text: "message from 444",
						},
					}, {
						Sender: struct {
							Id string "json:\"id\""
						}{
							Id: "666",
						},
						Recipient: struct {
							Id string "json:\"id\""
						}{
							Id: "123",
						},
						Timestamp: 1569262485349,
						Message: gometawebhooks.Message{
							Id:   "MESSAGE_ID",
							Text: "message from 666",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.MessagingHandler(testHandler{func() {
						scenario.trigger("messaging")
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"messaging": 3,
			},
		},
		{
			name:   "story mention",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object": "instagram",
				"entry": [
				  {
					"id": "123",
					"time": 1569262486134,
					"messaging": [
					  {
						"sender": {
						  "id": "567"
						},
						"recipient": {
						  "id": "123"
						},
						"timestamp": 1569262485349,
						"message": {
						  "mid": "890",
						  "attachments": [
							{
							  "type": "story_mention",
							  "payload": {
								"url": "<CDN_URL>"
							  }
							}
						  ]
						}
					  }
					]
				  }
				]
			  }`),
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []gometawebhooks.Messaging{{
						Sender: struct {
							Id string "json:\"id\""
						}{
							Id: "567",
						},
						Recipient: struct {
							Id string "json:\"id\""
						}{
							Id: "123",
						},
						Timestamp: 1569262485349,
						Message: gometawebhooks.Message{
							Id: "890",
							Attachments: []gometawebhooks.Attachment{{
								Type: "story_mention",
								Payload: struct {
									URL string "json:\"url,omitempty\""
								}{
									URL: "<CDN_URL>",
								},
							}},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.InstagramMessageHandler(testHandler{func() {
						scenario.trigger("message")
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"message": 1,
			},
		},
		{
			name:   "text message",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object": "instagram",
				"entry": [
				  {
					"id": "123",
					"time": 1569262486134,
					"messaging": [
					  {
						"sender": {
						  "id": "567"
						},
						"recipient": {
						  "id": "123"
						},
						"timestamp": 1569262485349,
						"message": {
						  "mid": "890",
						  "text": "Text in message"
						}
					  }
					]
				  }
				]
			  }`),
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []gometawebhooks.Messaging{{
						Sender: struct {
							Id string "json:\"id\""
						}{
							Id: "567",
						},
						Recipient: struct {
							Id string "json:\"id\""
						}{
							Id: "123",
						},
						Timestamp: 1569262485349,
						Message: gometawebhooks.Message{
							Id:   "890",
							Text: "Text in message",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.InstagramMessageHandler(testHandler{func() {
						scenario.trigger("message")
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"message": 1,
			},
		},
		{
			name:   "postback message",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object": "instagram",
				"entry": [
				  {
					"id": "123",
					"time": 1569262486134,
					"messaging": [
					  {
						"sender": {
						  "id": "567"
						},
						"recipient": {
						  "id": "123"
						},
						"timestamp": 1569262485349,
						"postback": {
							"mid":"MESSAGE-ID",           
							"title": "SELECTED-ICEBREAKER-REPLY-OR-CTA-BUTTON",
							"payload": "CUSTOMER-RESPONSE-PAYLOAD"
						}
					  }
					]
				  }
				]
			  }`),
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []gometawebhooks.Messaging{{
						Sender: struct {
							Id string "json:\"id\""
						}{
							Id: "567",
						},
						Recipient: struct {
							Id string "json:\"id\""
						}{
							Id: "123",
						},
						Timestamp: 1569262485349,
						Postback: gometawebhooks.Postback{
							Id:      "MESSAGE-ID",
							Title:   "SELECTED-ICEBREAKER-REPLY-OR-CTA-BUTTON",
							Payload: "CUSTOMER-RESPONSE-PAYLOAD",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.InstagramPostbackHandler(testHandler{func() {
						scenario.trigger("postback")
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"postback": 1,
			},
		},
		{
			name:   "referral message",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object": "instagram",
				"entry": [
				  {
					"id": "123",
					"time": 1569262486134,
					"messaging": [
					  {
						"sender": {
						  "id": "567"
						},
						"recipient": {
						  "id": "123"
						},
						"timestamp": 1569262485349,
						"referral": {
							"ref":"INFORMATION-INCLUDED-IN-REF-PARAMETER-OF-IGME-LINK",           
							"source": "IGME-SOURCE-LINK",
							"type": "OPEN_THREAD"
						}
					  }
					]
				  }
				]
			  }`),
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []gometawebhooks.Messaging{{
						Sender: struct {
							Id string "json:\"id\""
						}{
							Id: "567",
						},
						Recipient: struct {
							Id string "json:\"id\""
						}{
							Id: "123",
						},
						Timestamp: 1569262485349,
						Referral: gometawebhooks.Referral{
							Ref:    "INFORMATION-INCLUDED-IN-REF-PARAMETER-OF-IGME-LINK",
							Source: "IGME-SOURCE-LINK",
							Type:   "OPEN_THREAD",
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.InstagramReferralHandler(testHandler{func() {
						scenario.trigger("referral")
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"referral": 1,
			},
		},
		{
			name:   "quick reply message",
			method: http.MethodPost,
			body: strings.NewReader(`{
				"object": "instagram",
				"entry": [
				  {
					"id": "123",
					"time": 1569262486134,
					"messaging": [
					  {
						"sender": {
						  "id": "567"
						},
						"recipient": {
						  "id": "123"
						},
						"timestamp": 1569262485349,
						"message": {
						  "mid": "890",
						  "quick_reply": {
								"payload":"QR-PAYLOAD"
							}
						}
					  }
					]
				  }
				]
			  }`),
			expected: gometawebhooks.Event{
				Object: gometawebhooks.Instagram,
				Entry: []gometawebhooks.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []gometawebhooks.Messaging{{
						Sender: struct {
							Id string "json:\"id\""
						}{
							Id: "567",
						},
						Recipient: struct {
							Id string "json:\"id\""
						}{
							Id: "123",
						},
						Timestamp: 1569262485349,
						Message: gometawebhooks.Message{
							Id: "890",
							QuickReply: struct {
								Payload string "json:\"payload,omitempty\""
							}{"QR-PAYLOAD"},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []gometawebhooks.Option {
				return []gometawebhooks.Option{
					gometawebhooks.Options.InstagramMessageHandler(testHandler{func() {
						scenario.trigger("message")
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"message": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.test(t, func(t *testing.T) {
			hooks, req := scenario.setup(t)

			ctx := context.Background()

			result, err := hooks.Handle(ctx, req)

			scenario.assert(t, result, err)
		})
	}
}
