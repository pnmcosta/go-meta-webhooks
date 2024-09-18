package gometawebhooks

import (
	"context"
	"encoding/json"
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
		// one
		Id string `json:"mid,omitempty"`
		// or another
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
	Type    string `json:"type,omitempty"`
	Source  string `json:"source,omitempty"`
	Ref     string `json:"ref,omitempty"`
	Product struct {
		Id string `json:"id,omitempty"`
	} `product:"ref,omitempty"`
}

type Postback struct {
	Id       string   `json:"mid,omitempty"`
	Title    string   `json:"title,omitempty"`
	Payload  string   `json:"payload,omitempty"`
	Referral Referral `json:"referral,omitempty"`
}

type MessagingHeader struct {
	Sender struct {
		Id string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		Id string `json:"id"`
	} `json:"recipient"`
	Timestamp int64 `json:"timestamp"`
}

// https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messages
type MessagingMessage struct {
	MessagingHeader

	Message Message `json:"message,omitempty"`
}

// https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messaging-postbacks
type MessagingPostback struct {
	MessagingHeader

	Postback Postback `json:"postback,omitempty"`
}

// https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#igme
type MessagingReferral struct {
	MessagingHeader

	Referral Referral `json:"referral,omitempty"`
}

// @todo these
// https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#message-reactions
// https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messaging-seen

// Wrapper struct for types MessagingMessage, MessagingPostback and MessagingReferral
type Messaging struct {
	Type interface{} `json:"-"`
}

func (t *Messaging) UnmarshalJSON(b []byte) error {
	var message MessagingMessage
	if err := json.Unmarshal(b, &message); err == nil && message.Message.Id != "" {
		t.Type = message
		return nil
	}

	var postback MessagingPostback
	if err := json.Unmarshal(b, &postback); err == nil && postback.Postback.Id != "" {
		t.Type = postback
		return nil
	}

	var referral MessagingReferral
	if err := json.Unmarshal(b, &referral); err == nil && referral.Referral.Type != "" {
		t.Type = referral
		return nil
	}

	return ErrMessagingTypeNotImplemented
}

func (hooks Webhooks) messaging(ctx context.Context, object Object, entry Entry) error {
	if len(entry.Messaging) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(len(entry.Messaging))
	for _, messaging := range entry.Messaging {
		g.Go(func() error {
			return hooks.message(ctx, object, entry, messaging)
		})
	}
	return g.Wait()
}

func (h Webhooks) message(ctx context.Context, object Object, entry Entry, messaging Messaging) error {
	switch value := messaging.Type.(type) {
	case MessagingMessage:
		if h.ignoreEchoMessages && value.Message.IsEcho {
			return nil
		}

		if h.instagramMessageHandler == nil {
			return ErrInstagramMessageHandlerNotDefined
		}

		return h.instagramMessageHandler.InstagramMessage(ctx, object, entry, value)
	case MessagingPostback:
		if h.instagramPostbackHandler == nil {
			return ErrInstagramPostbackHandlerNotDefined
		}

		return h.instagramPostbackHandler.InstagramPostback(ctx, object, entry, value)
	case MessagingReferral:
		if h.instagramReferralHandler == nil {
			return ErrInstagramReferralHandlerNotDefined
		}

		return h.instagramReferralHandler.InstagramReferral(ctx, object, entry, value)
	default:
		// @note should not be hit cause Unmarshall ensures field is supported
		return ErrMessagingTypeNotImplemented
	}
}
