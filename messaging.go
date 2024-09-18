package gometawebhooks

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	ErrMessagingTypeNotImplemented        = errors.New("messaging type not implemented")
	ErrInstagramMessageHandlerNotDefined  = errors.New("instagram message handler not defined")
	ErrInstagramPostbackHandlerNotDefined = errors.New("instagram postback handler not defined")
	ErrInstagramReferralHandlerNotDefined = errors.New("instagram referral handler not defined")
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
	Messaging(context.Context, Object, Entry, Messaging) error
}

func (hooks Webhooks) messaging(ctx context.Context, object Object, entry Entry) error {
	if len(entry.Messaging) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(len(entry.Messaging))
	for _, messaging := range entry.Messaging {
		g.Go(func() error {
			return hooks.messagingHandler.Messaging(ctx, object, entry, messaging)
		})
	}
	return g.Wait()
}

func (h Webhooks) Messaging(ctx context.Context, object Object, entry Entry, messaging Messaging) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	default:
		return h.message(ctx, object, entry, messaging)
	}
}

func (h Webhooks) message(ctx context.Context, object Object, entry Entry, messaging Messaging) error {
	if h.messagingIgnoreEchos && messaging.Message.IsEcho {
		return nil
	}

	sent := unixTime(messaging.Timestamp)
	if messaging.Message.Id != "" {
		if h.instagramMessageHandler == nil {
			return ErrInstagramMessageHandlerNotDefined
		}

		return h.instagramMessageHandler.InstagramMessage(ctx, object, entry, messaging.Sender.Id, messaging.Recipient.Id, sent, messaging.Message)
	}

	if messaging.Postback.Id != "" {
		if h.instagramPostbackHandler == nil {
			return ErrInstagramPostbackHandlerNotDefined
		}

		return h.instagramPostbackHandler.InstagramPostback(ctx, object, entry, messaging.Sender.Id, messaging.Recipient.Id, sent, messaging.Postback)
	}

	if messaging.Referral.Type != "" {
		if h.instagramReferralHandler == nil {
			return ErrInstagramReferralHandlerNotDefined
		}

		return h.instagramReferralHandler.InstagramReferral(ctx, object, entry, messaging.Sender.Id, messaging.Recipient.Id, sent, messaging.Referral)
	}

	return ErrMessagingTypeNotImplemented
}
