package gometawebhooks

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrWebhooks       = errors.New("gometawebhooks")
	ErrApplyingOption = fmt.Errorf("error applying option: %w", ErrWebhooks)
)

// Webhooks instance contains all methods needed to process object events
type Webhooks struct {
	token  string
	secret string

	handleEntry func(ctx context.Context, object Object, entry Entry)

	handleInstagramChange       func(ctx context.Context, entry Entry, change Change)
	handleInstagramMention      func(ctx context.Context, entry Entry, mention MentionsFieldValue)
	handleInstagramStoryInsight func(ctx context.Context, entry Entry, storyInsights StoryInsightsFieldValue)

	handleInstagramMessaging func(ctx context.Context, entry Entry, messaging Messaging)
	handleInstagramMessage   func(ctx context.Context, sender string, recipient string, time int64, message Message)
	handleInstagramPostback  func(ctx context.Context, sender string, recipient string, time int64, postback Postback)
	handleInstagramReferral  func(ctx context.Context, sender string, recipient string, time int64, referral Referral)
}

// Creates and returns a webhooks instance
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
