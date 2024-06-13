package gometawebhooks

import (
	"context"
	"sync"
	"time"
)

type Message struct {
	Id          string       `json:"mid,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Referral    Referral     `json:"referral,omitempty"`

	IsDeleted     bool `json:"is_deleted,omitempty"`
	IsEcho        bool `json:"is_echo,omitempty"`
	IsUnsupported bool `json:"is_unsupported,omitempty"`

	ReplyTo struct {
		Id    string `json:"mid,omitempty"`
		Story struct {
			ID  string `json:"id,omitempty"`
			URL string `json:"url,omitempty"`
		} `json:"story,omitempty"`
	} `json:"reply_to,omitempty"`

	QuickReply struct {
		Payload string `json:"payload,omitempty"`
	} `json:"quick_reply,omitempty"`
}

type Attachment struct {
	Type    string            `json:"type,omitempty"`
	Payload AttachmentPayload `json:"payload,omitempty"`
}

type AttachmentPayload struct {
	URL         string `json:"url,omitempty"`
	Title       string `json:"title,omitempty"`
	StickerId   string `json:"sticker_id,omitempty"`
	ReelVideoId string `json:"reel_video_id,omitempty"`
}

type Referral struct {
	Type   string `json:"type,omitempty"`
	Source string `json:"source,omitempty"`
	Ref    string `json:"ref,omitempty"`
}

type Postback struct {
	Id       string   `json:"mid,omitempty"`
	Title    string   `json:"title,omitempty"`
	Payload  string   `json:"payload,omitempty"`
	Referral Referral `json:"referral,omitempty"`
}

type Messaging struct {
	Sender struct {
		Id string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		Id string `json:"id"`
	} `json:"recipient"`
	Timestamp int64    `json:"timestamp"`
	Message   Message  `json:"message,omitempty"`
	Postback  Postback `json:"postback,omitempty"`
	Referral  Referral `json:"referral,omitempty"`
}

type MessagingHandler interface {
	Messaging(ctx context.Context, object Object, entryId string, entryTime time.Time, messaging Messaging)
}

func (hooks Webhooks) messaging(ctx context.Context, object Object, entry Entry) {
	if len(entry.Messaging) == 0 {
		return
	}

	var wg sync.WaitGroup
out:
	for _, messaging := range entry.Messaging {
		select {
		case <-ctx.Done():
			break out
		default:
		}

		wg.Add(1)
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
