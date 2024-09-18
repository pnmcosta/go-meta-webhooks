package gometawebhooks

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrApplyingOption = errors.New("error applying option")
)

// Webhooks instance contains all methods needed to process object events
type Webhooks struct {
	token  string
	secret string

	headerSigName string

	entryHandler     EntryHandler
	changesHandler   ChangesHandler
	messagingHandler MessagingHandler

	instagramMessageHandler       InstagramMessageHandler
	instagramPostbackHandler      InstagramPostbackHandler
	instagramReferralHandler      InstagramReferralHandler
	instagramMentionHandler       InstagramMentionHandler
	instagramStoryInsightsHandler InstagramStoryInsightsHandler

	messagingIgnoreEchos bool
}

type WebhooksHandler interface {
	EntryHandler
	ChangesHandler
	MessagingHandler

	Handle(context.Context, Event) error
}

var _ WebhooksHandler = (*Webhooks)(nil)

// Creates and returns a webhooks instance
func New(options ...Option) (*Webhooks, error) {
	hooks := &Webhooks{}

	for _, opt := range options {
		if err := opt(hooks); err != nil {
			return nil, wrapErr(err, ErrApplyingOption)
		}
	}

	if hooks.headerSigName == "" {
		hooks.headerSigName = HeaderSignatureName
	}

	if hooks.entryHandler == nil {
		hooks.entryHandler = hooks
	}

	if hooks.changesHandler == nil {
		hooks.changesHandler = hooks
	}

	if hooks.messagingHandler == nil {
		hooks.messagingHandler = hooks
	}

	return hooks, nil
}

func wrapErr(err, target error) error {
	return fmt.Errorf("%w: %w", err, target)
}

func unixTime(timeMs int64) time.Time {
	return time.Unix(0, timeMs*int64(time.Millisecond))
}
