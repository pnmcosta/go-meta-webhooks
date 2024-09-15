package gometawebhooks

import (
	"context"
	"time"
)

type InstagramMessageHandler interface {
	InstagramMessage(ctx context.Context, sender, recipient string, sent time.Time, message Message) error
}

type InstagramPostbackHandler interface {
	InstagramPostback(ctx context.Context, sender, recipient string, sent time.Time, postback Postback) error
}

type InstagramReferralHandler interface {
	InstagramReferral(ctx context.Context, sender, recipient string, sent time.Time, referral Referral) error
}

type InstagramMentionHandler interface {
	InstagramMention(ctx context.Context, entryId string, entryTime time.Time, mention MentionsFieldValue) error
}

type InstagramStoryInsightsHandler interface {
	InstagramStoryInsights(ctx context.Context, entryId string, entryTime time.Time, storyInsights StoryInsightsFieldValue) error
}

type InstagramChangesHandler interface {
	InstagramMentionHandler
	InstagramStoryInsightsHandler
}

type InstagramMessagingHandler interface {
	InstagramMessageHandler
	InstagramPostbackHandler
	InstagramReferralHandler
}

type InstagramHandler interface {
	InstagramChangesHandler
	InstagramMessagingHandler
}
