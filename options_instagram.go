package gometawebhooks

import "context"

// Handle https://developers.facebook.com/docs/instagram-api/guides/mentions
func (MetaWebhookOptions) HandleInstagramMention(fn func(ctx context.Context, entry Entry, mention MentionsFieldValue)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramMention = fn
		return nil
	}
}

// Handle https://developers.facebook.com/docs/instagram-api/guides/webhooks#capturing-story-insights
func (MetaWebhookOptions) HandleInstagramStoryInsight(fn func(ctx context.Context, entry Entry, storyInsights StoryInsightsFieldValue)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramStoryInsight = fn
		return nil
	}
}

// Overrides this libraries default bevaviour to handle Instagram messaging echos, deleted, postbacks, referrals, reads and reactions
// Please note this will prevent other handle options from executing.
func (MetaWebhookOptions) HandleInstagramMessaging(fn func(ctx context.Context, entry Entry, messaging Messaging)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramMessaging = fn
		return nil
	}
}

// Handle https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messages
func (MetaWebhookOptions) HandleInstagramMessage(fn func(ctx context.Context, sender string, recipient string, time int64, message Message)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramMessage = fn
		return nil
	}
}

// Handle https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messaging-postbacks
func (MetaWebhookOptions) HandleInstagramPostback(fn func(ctx context.Context, sender string, recipient string, time int64, postback Postback)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramPostback = fn
		return nil
	}
}

// Handle https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#igme
func (MetaWebhookOptions) HandleInstagramReferral(fn func(ctx context.Context, sender string, recipient string, time int64, referral Referral)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramReferral = fn
		return nil
	}
}
