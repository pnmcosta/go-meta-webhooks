package gometawebhooks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChangesFieldNotImplemented              = errors.New("changes field not implemented")
	ErrInstagramMentionHandlerNotDefined       = errors.New("instagram mentions handler not defined")
	ErrInstagramStoryInsightsHandlerNotDefined = errors.New("instagram story insights handler not defined")
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
	Changes(context.Context, Object, Entry, Change) error
}

func (c *Change) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if fieldRaw, ok := raw["field"]; ok {
		if err := json.Unmarshal(fieldRaw, &c.Field); err != nil {
			return err
		}
	}

	if valueRaw, ok := raw["value"]; ok {
		switch c.Field {
		case "mentions":
			var value MentionsFieldValue
			if err := json.Unmarshal(valueRaw, &value); err != nil {
				return err
			}
			c.Value = value
		case "story_insights":
			var value StoryInsightsFieldValue
			if err := json.Unmarshal(valueRaw, &value); err != nil {
				return err
			}
			c.Value = value
		default:
			return fmt.Errorf("'%s': %w", c.Field, ErrChangesFieldNotImplemented)
		}
	}

	return nil
}

func (hooks Webhooks) changes(ctx context.Context, object Object, entry Entry) error {
	if len(entry.Changes) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(len(entry.Changes))
	for _, change := range entry.Changes {
		g.Go(func() error {
			return hooks.changesHandler.Changes(ctx, object, entry, change)
		})
	}

	return g.Wait()
}

func (h Webhooks) Changes(ctx context.Context, object Object, entry Entry, change Change) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	default:
		return h.change(ctx, object, entry, change)
	}
}
func (h Webhooks) change(ctx context.Context, object Object, entry Entry, change Change) error {
	switch value := change.Value.(type) {
	case MentionsFieldValue:
		if h.instagramMentionHandler == nil {
			return ErrInstagramMentionHandlerNotDefined
		}
		return h.instagramMentionHandler.InstagramMention(ctx, object, entry, value)
	case StoryInsightsFieldValue:
		if h.instagramStoryInsightsHandler == nil {
			return ErrInstagramStoryInsightsHandlerNotDefined
		}
		return h.instagramStoryInsightsHandler.InstagramStoryInsights(ctx, object, entry, value)
	default:
		// @note should not be hit cause Unmarshall ensures field is supported
		return fmt.Errorf("'%s': %w", change.Field, ErrChangesFieldNotImplemented)
	}
}
