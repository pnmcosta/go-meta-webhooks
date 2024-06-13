package gometawebhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

var (
	ErrParsingChanges           = fmt.Errorf("error parsing changes payload: %w", ErrWebhooks)
	ErrChangesFieldNotSupported = fmt.Errorf("field not supported: %w", ErrWebhooks)
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
	Changes(ctx context.Context, object Object, entry Entry, change Change)
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

func (hooks Webhooks) changes(ctx context.Context, object Object, entry Entry) {
	if len(entry.Changes) == 0 {
		return
	}

	var wg sync.WaitGroup

out:
	for _, change := range entry.Changes {
		select {
		case <-ctx.Done():
			break out
		default:
		}
		wg.Add(1)

		change := change

		go func() {
			defer wg.Done()

			hooks.changesHandler.Changes(ctx, object, entry, change)
		}()
	}

	wg.Wait()
}

func (h defaultHandler) Changes(ctx context.Context, object Object, entry Entry, change Change) {
	if object != Instagram {
		return
	}

	sent := unixTime(entry.Time)

	switch value := change.Value.(type) {
	case MentionsFieldValue:
		if h.hooks.instagramMentionHandler == nil {
			return
		}
		h.hooks.instagramMentionHandler.InstagramMention(ctx, entry.Id, sent, value)
	case StoryInsightsFieldValue:
		if h.hooks.instagramStoryInsightsHandler == nil {
			return
		}
		h.hooks.instagramStoryInsightsHandler.InstagramStoryInsights(ctx, entry.Id, sent, value)
	}
}
