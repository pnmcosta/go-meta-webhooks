package gometawebhooks

import "context"

func (MetaWebhookOptions) OnInstagramMention(fn func(context.Context, Entry, MentionsFieldValue)) Option {
	return func(hook *Webhooks) error {
		hook.onInstagramMention = fn
		return nil
	}
}

func (MetaWebhookOptions) OnInstagramStoryInsight(fn func(context.Context, Entry, StoryInsightsFieldValue)) Option {
	return func(hook *Webhooks) error {
		hook.onInstagramStoryInsight = fn
		return nil
	}
}

func (MetaWebhookOptions) OnInstagramMessaging(fn func(context.Context, Entry, Messaging)) Option {
	return func(hook *Webhooks) error {
		hook.onInstagramMessaging = fn
		return nil
	}
}

func (MetaWebhookOptions) OnInstagramMessage(fn func(context.Context, string, string, int64, Message)) Option {
	return func(hook *Webhooks) error {
		hook.onInstagramMessage = fn
		return nil
	}
}

func (MetaWebhookOptions) OnInstagramPostback(fn func(context.Context, string, string, int64, Postback)) Option {
	return func(hook *Webhooks) error {
		hook.onInstagramPostback = fn
		return nil
	}
}

func (MetaWebhookOptions) OnInstagramReferral(fn func(context.Context, string, string, int64, Referral)) Option {
	return func(hook *Webhooks) error {
		hook.onInstagramReferral = fn
		return nil
	}
}
