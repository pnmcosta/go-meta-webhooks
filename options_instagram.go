package gometawebhooks

// Sets the InstagramMentionHandler, see https://developers.facebook.com/docs/instagram-api/guides/mentions
func (MetaWebhookOptions) InstagramMentionHandler(fn InstagramMentionHandler) Option {
	return func(hooks *Webhooks) error {
		hooks.instagramMentionHandler = fn
		return nil
	}
}

// Sets the InstagramStoryInsightsHandler, see https://developers.facebook.com/docs/instagram-api/guides/webhooks#capturing-story-insights
func (MetaWebhookOptions) InstagramStoryInsightsHandler(fn InstagramStoryInsightsHandler) Option {
	return func(hooks *Webhooks) error {
		hooks.instagramStoryInsightsHandler = fn
		return nil
	}
}

// Sets the InstagramMessageHandler, see https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messages
func (MetaWebhookOptions) InstagramMessageHandler(fn InstagramMessageHandler) Option {
	return func(hooks *Webhooks) error {
		hooks.instagramMessageHandler = fn
		return nil
	}
}

// Sets the InstagramPostbackHandler, see https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#messaging-postbacks
func (MetaWebhookOptions) InstagramPostbackHandler(fn InstagramPostbackHandler) Option {
	return func(hooks *Webhooks) error {
		hooks.instagramPostbackHandler = fn
		return nil
	}
}

// Sets the InstagramReferralHandler, see https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook/#igme
func (MetaWebhookOptions) InstagramReferralHandler(fn InstagramReferralHandler) Option {
	return func(hooks *Webhooks) error {
		hooks.instagramReferralHandler = fn
		return nil
	}
}

// Sets all InstagramMessaging handlers
func (MetaWebhookOptions) InstagramMessagingHandler(fn InstagramMessagingHandler) Option {
	return func(hooks *Webhooks) error {
		hooks.instagramMessageHandler = fn
		hooks.instagramPostbackHandler = fn
		hooks.instagramReferralHandler = fn
		return nil
	}
}

// Sets all InstagramChanges handlers
func (MetaWebhookOptions) InstagramChangesHandler(fn InstagramChangesHandler) Option {
	return func(hooks *Webhooks) error {
		hooks.instagramMentionHandler = fn
		hooks.instagramStoryInsightsHandler = fn
		return nil
	}
}

// Sets all Instagram handlers
func (MetaWebhookOptions) InstagramHandler(fn InstagramHandler) Option {
	return func(hooks *Webhooks) error {
		hooks.instagramMessageHandler = fn
		hooks.instagramPostbackHandler = fn
		hooks.instagramReferralHandler = fn
		hooks.instagramMentionHandler = fn
		hooks.instagramStoryInsightsHandler = fn
		return nil
	}
}
