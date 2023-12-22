package gometawebhooks

// Option is a configuration option for the webhook
type Option func(*Webhooks) error

// Options is a namespace var for configuration options
var Options = MetaWebhookOptions{}

// MetaWebhookOptions is a namespace for configuration option methods
type MetaWebhookOptions struct{}

// Sets the Facebook APP Secret
func (MetaWebhookOptions) Secret(secret string) Option {
	return func(hook *Webhooks) error {
		hook.secret = secret
		return nil
	}
}

// Sets the Facebook APP webhook subscription verify token
func (MetaWebhookOptions) Token(token string) Option {
	return func(hook *Webhooks) error {
		hook.token = token
		return nil
	}
}

// Overrides the default EntryHandler, please note this will override object handler options.
func (MetaWebhookOptions) EntryHandler(h EntryHandler) Option {
	return func(hook *Webhooks) error {
		hook.entryHandler = h
		return nil
	}
}

// Overrides the default ChangesHandler, please note this will override object handler options.
func (MetaWebhookOptions) ChangesHandler(h ChangesHandler) Option {
	return func(hook *Webhooks) error {
		hook.changesHandler = h
		return nil
	}
}

// Overrides the default MessagingHandler, please note this will override object handler options.
func (MetaWebhookOptions) MessagingHandler(h MessagingHandler) Option {
	return func(hook *Webhooks) error {
		hook.messagingHandler = h
		return nil
	}
}
