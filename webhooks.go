package gometawebhooks

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrWebhooks       = errors.New("GoMetaWebhooks")
	ErrApplyingOption = fmt.Errorf("error applying option: %w", ErrWebhooks)
)

// Webhooks instance contains all methods needed to process object events
type Webhooks struct {
	token  string
	secret string

	onChange    func(context.Context, Object, Entry, Change)
	onMessaging func(context.Context, Object, Entry, Messaging)

	onInstagramMention      func(context.Context, Entry, MentionsFieldValue)
	onInstagramStoryInsight func(context.Context, Entry, StoryInsightsFieldValue)
	onInstagramMessaging    func(context.Context, Entry, Messaging)

	onInstagramMessage  func(context.Context, string, string, int64, Message)
	onInstagramPostback func(context.Context, string, string, int64, Postback)
	onInstagramReferral func(context.Context, string, string, int64, Referral)
}

// Creates and returns a Webhooks instance
func NewWebhooks(options ...Option) (*Webhooks, error) {
	hooks := &Webhooks{}

	for _, opt := range options {
		if err := opt(hooks); err != nil {
			return nil, wrapErr(err, ErrApplyingOption)
		}
	}

	if err := hooks.compileSchema(); err != nil {
		return nil, err
	}

	return hooks, nil
}

func wrapErr(err, target error) error {
	return fmt.Errorf("%s: %w", err, target)
}
