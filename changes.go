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
	ErrChangesFieldNotSupported = fmt.Errorf("field not supported: %w", ErrWebhooks)
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
		}
	}

	return nil
}

func (hook Webhooks) changes(ctx context.Context, object Object, entry Entry) {
	if len(entry.Changes) == 0 {
		return
	}

	instagramChange := hook.handleInstagramChange
	if instagramChange == nil {
		instagramChange = hook.handleInstagramChangeDefault
	}

	var wg sync.WaitGroup
	wg.Add(len(entry.Changes))

	for _, change := range entry.Changes {
		select {
		case <-ctx.Done():
			break
		default:
		}

		change := change

		go func() {
			defer wg.Done()

			instagramChange(ctx, entry, change)
		}()
	}

	wg.Wait()
}

func (hook Webhooks) handleInstagramChangeDefault(ctx context.Context, entry Entry, change Change) {
	switch value := change.Value.(type) {
	case MentionsFieldValue:
		if hook.handleInstagramMention != nil {
			hook.handleInstagramMention(ctx, entry, value)
		}
	case StoryInsightsFieldValue:
		if hook.handleInstagramStoryInsight != nil {
			hook.handleInstagramStoryInsight(ctx, entry, value)
		}
	default:
		log.Printf("meta webhook event %v entry %s change field %s not supported\n", Instagram, entry.Id, change.Field)
	}
}
