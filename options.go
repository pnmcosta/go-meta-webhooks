package gometawebhooks

// Option is a configuration option for the webhook
type Option func(*Webhooks) error

// Options is a namespace var for configuration options
var Options = MetaWebhookOptions{}

// MetaWebhookOptions is a namespace for configuration option methods
type MetaWebhookOptions struct{}

// Sets the Facebook APP Secret
func (MetaWebhookOptions) Secret(secret string) Option {
	return func(hooks *Webhooks) error {
		hooks.secret = secret
		return nil
	}
}

// Sets the Facebook APP webhook subscription verify token
func (MetaWebhookOptions) Token(token string) Option {
	return func(hooks *Webhooks) error {
		hooks.token = token
		return nil
	}
}

func (MetaWebhookOptions) IgnoreEchoMessages(ignore bool) Option {
	return func(hooks *Webhooks) error {
		hooks.ignoreEchoMessages = ignore
		return nil
	}
}

// Ensures embedded JSON schema is compiled
func (MetaWebhookOptions) CompileSchema() Option {
	return func(hooks *Webhooks) error {
		if err := hooks.compileSchema(); err != nil {
			return err
		}
		return nil
	}
}

// Sets a custom header signature name
func (MetaWebhookOptions) CustomHeaderSigName(name string) Option {
	return func(hooks *Webhooks) error {
		hooks.headerSigName = name
		return nil
	}
}
