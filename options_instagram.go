package gometawebhooks

// Handle https://developers.facebook.com/docs/instagram-api/guides/mentions
func (MetaWebhookOptions) InstagramMentionHandler(fn InstagramMentionHandler) Option {
	return func(hook *Webhooks) error {
		hook.instagramMentionHandler = fn
		return nil
	}
}

// Handle https://developers.facebook.com/docs/instagram-api/guides/webhooks#capturing-story-insights
func (MetaWebhookOptions) InstagramStoryInsightsHandler(fn InstagramStoryInsightsHandler) Option {
	return func(hook *Webhooks) error {
		hook.instagramStoryInsightsHandler = fn
		return nil
	}
}

// Handle https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messages
func (MetaWebhookOptions) InstagramMessageHandler(fn InstagramMessageHandler) Option {
	return func(hook *Webhooks) error {
		hook.instagramMessageHandler = fn
		return nil
	}
}

// Handle https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messaging-postbacks
func (MetaWebhookOptions) InstagramPostbackHandler(fn InstagramPostbackHandler) Option {
	return func(hook *Webhooks) error {
		hook.instagramPostbackHandler = fn
		return nil
	}
}

// Handle https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#igme
func (MetaWebhookOptions) InstagramReferralHandler(fn InstagramReferralHandler) Option {
	return func(hook *Webhooks) error {
		hook.instagramReferralHandler = fn
		return nil
	}
}

func (MetaWebhookOptions) InstagramMessagingHandler(fn InstagramMessagingHandler) Option {
	return func(hook *Webhooks) error {
		hook.instagramMessageHandler = fn
		hook.instagramPostbackHandler = fn
		hook.instagramReferralHandler = fn
		return nil
	}
}

func (MetaWebhookOptions) InstagramChangesHandler(fn InstagramChangesHandler) Option {
	return func(hook *Webhooks) error {
		hook.instagramMentionHandler = fn
		hook.instagramStoryInsightsHandler = fn
		return nil
	}
}

func (MetaWebhookOptions) InstagramHandler(fn InstagramHandler) Option {
	return func(hook *Webhooks) error {
		hook.instagramMessageHandler = fn
		hook.instagramPostbackHandler = fn
		hook.instagramReferralHandler = fn
		hook.instagramMentionHandler = fn
		hook.instagramStoryInsightsHandler = fn
		return nil
	}
}
