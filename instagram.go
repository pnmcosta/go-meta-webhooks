package gometawebhooks

import (
	"context"
)

type InstagramMessageHandler interface {
	InstagramMessage(ctx context.Context, object Object, entry Entry, message MessagingMessage) error
}

type InstagramPostbackHandler interface {
	InstagramPostback(ctx context.Context, object Object, entry Entry, postback MessagingPostback) error
}

type InstagramReferralHandler interface {
	InstagramReferral(ctx context.Context, object Object, entry Entry, referral MessagingReferral) error
}

type InstagramMentionHandler interface {
	InstagramMention(ctx context.Context, object Object, entry Entry, mention Mention) error
}

type InstagramStoryInsightsHandler interface {
	InstagramStoryInsights(ctx context.Context, object Object, entry Entry, storyInsights StoryInsights) error
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
