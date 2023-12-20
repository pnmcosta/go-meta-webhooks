package gometawebhooks

import (
	"context"
	"log"
	"sync"
)

type Message struct {
	Id          string       `json:"mid"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Referral    Referral     `json:"referral,omitempty"`

	IsDeleted     bool `json:"is_deleted,omitempty"`
	IsEcho        bool `json:"is_echo,omitempty"`
	IsUnsupported bool `json:"is_unsupported,omitempty"`

	ReplyTo struct {
		Id    string `json:"mid,omitempty"`
		Story struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		} `json:"story,omitempty"`
	} `json:"reply_to,omitempty"`
}

type Attachment struct {
	Type    string `json:"type"`
	Payload struct {
		URL string `json:"url"`
	} `json:"payload"`
}

type Referral struct {
	Type   string `json:"type"`
	Source string `json:"source"`
	Ref    string `json:"ref"`
}

type Postback struct {
	Id       string   `json:"mid"`
	Title    string   `json:"title"`
	Payload  string   `json:"payload"`
	Referral Referral `json:"referral"`
}

type Messaging struct {
	Sender struct {
		Id string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		Id string `json:"id"`
	} `json:"recipient"`
	Timestamp int64    `json:"timestamp"`
	Message   Message  `json:"message"`
	Postback  Postback `json:"postback"`
	Referral  Referral `json:"referral"`
}

func (hook Webhooks) messaging(ctx context.Context, object Object, entry Entry) {
	if len(entry.Messaging) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, messaging := range entry.Messaging {
		messaging := messaging
		wg.Add(1)

		go func(messaging Messaging) {
			defer wg.Done()

			fn := hook.onMessaging
			if fn == nil {
				fn = hook.defaultOnMessaging
			}

			fn(ctx, object, entry, messaging)
		}(messaging)

		wg.Wait()
	}
}

func (hook Webhooks) defaultOnMessaging(ctx context.Context, object Object, entry Entry, messaging Messaging) {
	switch object {
	case Instagram:
		fn := hook.onInstagramMessaging
		if fn == nil {
			fn = hook.defaultOnInstagramMessaging
		}

		fn(ctx, entry, messaging)
	default:
		log.Printf("meta webhook event object %s messaging not supported\n", object)
	}
}
