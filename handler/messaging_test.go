package handler_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/pnmcosta/go-meta-webhooks/handler"
)

func TestHandleMessaging(t *testing.T) {
	t.Parallel()
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
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []handler.Messaging{
						{
							Type: handler.MessagingMessage{
								MessagingHeader: handler.MessagingHeader{
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
								},
								Message: handler.Message{
									Id:   "MESSAGE_ID",
									Text: "message from 567",
								},
							},
						},
						{
							Type: handler.MessagingMessage{
								MessagingHeader: handler.MessagingHeader{
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
								},
								Message: handler.Message{
									Id:   "MESSAGE_ID",
									Text: "message from 444",
								},
							},
						},
						{
							Type: handler.MessagingMessage{
								MessagingHeader: handler.MessagingHeader{
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
								},
								Message: handler.Message{
									Id:   "MESSAGE_ID",
									Text: "message from 666",
								},
							},
						}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramMessageHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("message")
						return nil
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"message": 3,
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
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []handler.Messaging{{
						Type: handler.MessagingMessage{
							MessagingHeader: handler.MessagingHeader{
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
							},
							Message: handler.Message{
								Id: "890",
								Attachments: []handler.Attachment{{
									Type: "story_mention",
									Payload: handler.AttachmentPayload{
										URL: "<CDN_URL>",
									},
								}},
							},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramMessageHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("message")
						return nil
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
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []handler.Messaging{{
						Type: handler.MessagingMessage{
							MessagingHeader: handler.MessagingHeader{
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
							},
							Message: handler.Message{
								Id:   "890",
								Text: "Text in message",
							},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramMessageHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("message")
						return nil
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
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []handler.Messaging{{
						Type: handler.MessagingPostback{
							MessagingHeader: handler.MessagingHeader{
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
							},
							Postback: handler.Postback{
								Id:      "MESSAGE-ID",
								Title:   "SELECTED-ICEBREAKER-REPLY-OR-CTA-BUTTON",
								Payload: "CUSTOMER-RESPONSE-PAYLOAD",
							},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramPostbackHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("postback")
						return nil
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
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []handler.Messaging{{
						Type: handler.MessagingReferral{
							MessagingHeader: handler.MessagingHeader{
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
							},
							Referral: handler.Referral{
								Ref:    "INFORMATION-INCLUDED-IN-REF-PARAMETER-OF-IGME-LINK",
								Source: "IGME-SOURCE-LINK",
								Type:   "OPEN_THREAD",
							},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramReferralHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("referral")
						return nil
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
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []handler.Messaging{{
						Type: handler.MessagingMessage{
							MessagingHeader: handler.MessagingHeader{
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
							},
							Message: handler.Message{
								Id: "890",
								QuickReply: struct {
									Payload string "json:\"payload,omitempty\""
								}{"QR-PAYLOAD"},
							},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramMessageHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("message")
						return nil
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"message": 1,
			},
		},
		{
			name:   "reel message attachment",
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
							  "type": "reel",
							  "payload": {
								"url": "<CDN_URL>",
								"title": "reel title",
								"reel_video_id": "123"
							  }
							}
						  ]
						}
					  }
					]
				  }
				]
			  }`),
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []handler.Messaging{{
						Type: handler.MessagingMessage{
							MessagingHeader: handler.MessagingHeader{
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
							},
							Message: handler.Message{
								Id: "890",
								Attachments: []handler.Attachment{{
									Type: "reel",
									Payload: handler.AttachmentPayload{
										URL:         "<CDN_URL>",
										Title:       "reel title",
										ReelVideoId: "123",
									},
								}},
							},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramMessageHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("message")
						return nil
					}}),
				}
			},
			expectedHandlers: map[string]int{
				"message": 1,
			},
		},
		{
			name:   "ig_reel message attachment",
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
							  "type": "ig_reel",
							  "payload": {
								"url": "<CDN_URL>",
								"title": "reel title",
								"reel_video_id": "123"
							  }
							}
						  ]
						}
					  }
					]
				  }
				]
			  }`),
			expected: handler.Event{
				Object: handler.Instagram,
				Entry: []handler.Entry{{
					Id:   "123",
					Time: 1569262486134,
					Messaging: []handler.Messaging{{
						Type: handler.MessagingMessage{
							MessagingHeader: handler.MessagingHeader{
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
							},
							Message: handler.Message{
								Id: "890",
								Attachments: []handler.Attachment{{
									Type: "ig_reel",
									Payload: handler.AttachmentPayload{
										URL:         "<CDN_URL>",
										Title:       "reel title",
										ReelVideoId: "123",
									},
								}},
							},
						},
					}},
				}},
			},
			options: func(scenario *hookScenario) []handler.Option {
				return []handler.Option{
					handler.Options.CompileSchema(),
					handler.Options.InstagramMessageHandler(testHandler{func(ctx context.Context) error {
						scenario.trigger("message")
						return nil
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

			result, payload, err := hooks.HandleRequest(ctx, req)

			scenario.assert(t, result, payload, err)
		})
	}
}
