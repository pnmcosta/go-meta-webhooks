package gometawebhooks

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var (
	ErrParsingChanges                          = fmt.Errorf("parsing changes payload: %w", ErrWebhooks)
	ErrChangesFieldNotSupported                = fmt.Errorf("changes field not supported: %w", ErrWebhooks)
	ErrChangesTypeNotImplemented               = fmt.Errorf("changes type not implemented: %w", ErrWebhooks)
	ErrInstagramMentionHandlerNotDefined       = fmt.Errorf("instagram mentions handler not defined: %w", ErrWebhooks)
	ErrInstagramStoryInsightsHandlerNotDefined = fmt.Errorf("instagram story insights handler not defined: %w", ErrWebhooks)
)

type Change struct {
	Field string      `json:"field,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type MentionsFieldValue struct {
	MediaID   string `json:"media_id"`
	CommentID string `json:"comment_id"`
}

type StoryInsightsFieldValue struct {
	MediaID     string `json:"media_id"`
	Exits       int    `json:"exits"`
	Replies     int    `json:"replies"`
	Reach       int    `json:"reach"`
	TapsForward int    `json:"taps_forward"`
	TapsBack    int    `json:"taps_back"`
	Impressions int    `json:"impressions"`
}

type ChangesHandler interface {
	Changes(ctx context.Context, object Object, entry Entry, change Change) error
}

func (c *Change) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return wrapErr(err, ErrParsingChanges)
	}

	if fieldRaw, ok := raw["field"]; ok {
		if err := json.Unmarshal(fieldRaw, &c.Field); err != nil {
			return wrapErr(err, ErrParsingChanges)
		}
	}

	if valueRaw, ok := raw["value"]; ok {
		switch c.Field {
		case "mentions":
			var value MentionsFieldValue
			if err := json.Unmarshal(valueRaw, &value); err != nil {
				return wrapErr(err, ErrParsingChanges)
			}
			c.Value = value
		case "story_insights":
			var value StoryInsightsFieldValue
			if err := json.Unmarshal(valueRaw, &value); err != nil {
				return wrapErr(err, ErrParsingChanges)
			}
			c.Value = value
		}
	}

	return nil
}

func (hooks Webhooks) changes(ctx context.Context, object Object, entry Entry) error {
	if len(entry.Changes) == 0 {
		return nil
	}

	g := new(errgroup.Group)
	g.SetLimit(1)
	for _, change := range entry.Changes {
		g.Go(func() error {
			return hooks.changesHandler.Changes(ctx, object, entry, change)
		})
	}

	return g.Wait()
}

func (h Webhooks) Changes(ctx context.Context, object Object, entry Entry, change Change) error {
	if object != Instagram {
		return fmt.Errorf("'%s': %w", object, ErrChangesFieldNotSupported)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return h.change(ctx, entry, change)
	}
}
func (h Webhooks) change(ctx context.Context, entry Entry, change Change) error {
	sent := unixTime(entry.Time)

	switch value := change.Value.(type) {
	case MentionsFieldValue:
		if h.instagramMentionHandler == nil {
			return ErrInstagramMentionHandlerNotDefined
		}
		return h.instagramMentionHandler.InstagramMention(ctx, entry.Id, sent, value)
	case StoryInsightsFieldValue:
		if h.instagramStoryInsightsHandler == nil {
			return ErrInstagramStoryInsightsHandlerNotDefined
		}
		return h.instagramStoryInsightsHandler.InstagramStoryInsights(ctx, entry.Id, sent, value)
	default:
		return fmt.Errorf("'%s': %w", change.Field, ErrChangesTypeNotImplemented)
	}
}
