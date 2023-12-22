package gometawebhooks

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrWebhooks       = errors.New("gometawebhooks")
	ErrApplyingOption = fmt.Errorf("error applying option: %w", ErrWebhooks)
)

// Webhooks instance contains all methods needed to process object events
type Webhooks struct {
	token  string
	secret string

	entryHandler     EntryHandler
	changesHandler   ChangesHandler
	messagingHandler MessagingHandler

	instagramMessageHandler       InstagramMessageHandler
	instagramPostbackHandler      InstagramPostbackHandler
	instagramReferralHandler      InstagramReferralHandler
	instagramMentionHandler       InstagramMentionHandler
	instagramStoryInsightsHandler InstagramStoryInsightsHandler
}

// Creates and returns a webhooks instance
func New(options ...Option) (*Webhooks, error) {
	hooks := &Webhooks{}

	for _, opt := range options {
		if err := opt(hooks); err != nil {
			return nil, wrapErr(err, ErrApplyingOption)
		}
	}

	if err := hooks.compileSchema(); err != nil {
		return nil, err
	}

	handler := defaultHandler{hooks}
	if hooks.entryHandler == nil {
		hooks.entryHandler = handler
	}

	if hooks.changesHandler == nil {
		hooks.changesHandler = handler
	}

	if hooks.messagingHandler == nil {
		hooks.messagingHandler = handler
	}

	return hooks, nil
}

func wrapErr(err, target error) error {
	return fmt.Errorf("%s: %w", err, target)
}

func unixTime(timeMs int64) time.Time {
	return time.Unix(0, timeMs*int64(time.Millisecond))
}
