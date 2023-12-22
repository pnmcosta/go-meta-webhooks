package gometawebhooks

import (
	"context"
	"time"
)

type InstagramMentionHandler interface {
	InstagramMention(ctx context.Context, entryId string, entryTime time.Time, mention MentionsFieldValue)
}

type InstagramStoryInsightsHandler interface {
	InstagramStoryInsights(ctx context.Context, entryId string, entryTime time.Time, storyInsights StoryInsightsFieldValue)
}

type InstagramChangesHandler interface {
	InstagramMentionHandler
	InstagramStoryInsightsHandler
}
