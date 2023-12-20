package gometawebhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

var (
	ErrParsingChanges           = fmt.Errorf("error parsing changes payload: %w", ErrWebhooks)
	ErrChangesFieldNotSupported = fmt.Errorf("changes field not supported: %w", ErrWebhooks)
)

type Change struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
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
		default:
			return ErrChangesFieldNotSupported
		}
	}

	return nil
}

func (hook Webhooks) changes(ctx context.Context, object Object, entry Entry) {
	if len(entry.Changes) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, change := range entry.Changes {
		change := change
		wg.Add(1)

		go func(change Change) {
			defer wg.Done()

			fn := hook.onChange
			if fn == nil {
				fn = hook.defaultOnChange
			}

			fn(ctx, object, entry, change)
		}(change)

		wg.Wait()
	}
}

func (hook Webhooks) defaultOnChange(ctx context.Context, object Object, entry Entry, change Change) {
	switch object {
	case Instagram:
		switch value := change.Value.(type) {
		case MentionsFieldValue:
			if hook.onInstagramMention != nil {
				hook.onInstagramMention(ctx, entry, value)
			}
		case StoryInsightsFieldValue:
			if hook.onInstagramStoryInsight != nil {
				hook.onInstagramStoryInsight(ctx, entry, value)
			}
		default:
			log.Printf("meta webhook event instagram field %s change not supported\n", change.Field)
		}
	default:
		log.Printf("meta webhook event object %s change not supported\n", object)
	}
}
