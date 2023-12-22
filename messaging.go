package gometawebhooks

import (
	"context"
	"sync"
	"time"
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

type MessagingHandler interface {
	Messaging(ctx context.Context, object Object, entryId string, entryTime time.Time, messaging Messaging)
}

func (hooks Webhooks) messaging(ctx context.Context, object Object, entry Entry) {
	if len(entry.Messaging) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(entry.Messaging))
	for _, messaging := range entry.Messaging {
		select {
		case <-ctx.Done():
			break
		default:
		}

		messaging := messaging
		go func() {
			defer wg.Done()

			hooks.messagingHandler.Messaging(ctx, object, entry.Id, unixTime(entry.Time), messaging)
		}()
	}

	wg.Wait()
}

func (h defaultHandler) Messaging(ctx context.Context, object Object, entryId string, entryTime time.Time, messaging Messaging) {
	if object != Instagram {
		return
	}

	sent := unixTime(messaging.Timestamp)
	if messaging.Message.Id != "" {
		if messaging.Message.IsEcho {
			return
		}

		if h.hooks.instagramMessageHandler != nil {
			h.hooks.instagramMessageHandler.InstagramMessage(ctx, messaging.Sender.Id, messaging.Recipient.Id, sent, messaging.Message)
		}

		return
	}

	if messaging.Postback.Id != "" {
		if h.hooks.instagramPostbackHandler != nil {
			h.hooks.instagramPostbackHandler.InstagramPostback(ctx, messaging.Sender.Id, messaging.Recipient.Id, sent, messaging.Postback)
		}
		return
	}

	if messaging.Referral.Type != "" {
		if h.hooks.instagramReferralHandler != nil {
			h.hooks.instagramReferralHandler.InstagramReferral(ctx, messaging.Sender.Id, messaging.Recipient.Id, sent, messaging.Referral)
		}
		return
	}
}
