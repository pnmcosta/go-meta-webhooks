package gometawebhooks

import "context"

func (MetaWebhookOptions) HandleInstagramMention(fn func(ctx context.Context, entry Entry, mention MentionsFieldValue)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramMention = fn
		return nil
	}
}

func (MetaWebhookOptions) HandleInstagramStoryInsight(fn func(ctx context.Context, entry Entry, storyInsights StoryInsightsFieldValue)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramStoryInsight = fn
		return nil
	}
}

func (MetaWebhookOptions) HandleInstagramMessaging(fn func(ctx context.Context, entry Entry, messaging Messaging)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramMessaging = fn
		return nil
	}
}

func (MetaWebhookOptions) HandleInstagramMessage(fn func(ctx context.Context, sender string, recipient string, time int64, message Message)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramMessage = fn
		return nil
	}
}

func (MetaWebhookOptions) HandleInstagramPostback(fn func(ctx context.Context, sender string, recipient string, time int64, postback Postback)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramPostback = fn
		return nil
	}
}

func (MetaWebhookOptions) HandleInstagramReferral(fn func(ctx context.Context, sender string, recipient string, time int64, referral Referral)) Option {
	return func(hook *Webhooks) error {
		hook.handleInstagramReferral = fn
		return nil
	}
}
