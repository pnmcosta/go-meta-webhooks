package gometawebhooks

import "context"

// Option is a configuration option for the webhook
type Option func(*Webhooks) error

// Options is a namespace var for configuration options
var Options = MetaWebhookOptions{}

// MetaWebhookOptions is a namespace for configuration option methods
type MetaWebhookOptions struct{}

// Secret registers the Facebook APP Secret
func (MetaWebhookOptions) Secret(secret string) Option {
	return func(hook *Webhooks) error {
		hook.secret = secret
		return nil
	}
}

// Token registers the Facebook verify_token
func (MetaWebhookOptions) Token(token string) Option {
	return func(hook *Webhooks) error {
		hook.token = token
		return nil
	}
}

func (MetaWebhookOptions) HandleChange(fn func(context.Context, Object, Entry, Change)) Option {
	return func(hook *Webhooks) error {
		hook.handleChange = fn
		return nil
	}
}

func (MetaWebhookOptions) HandleMessaging(fn func(context.Context, Object, Entry, Messaging)) Option {
	return func(hook *Webhooks) error {
		hook.handleMessaging = fn
		return nil
	}
}
