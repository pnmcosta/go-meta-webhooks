package gometawebhooks

import "context"

func (hook Webhooks) handleInstagramMessagingDefault(ctx context.Context, entry Entry, messaging Messaging) {
	if messaging.Message.IsEcho {
		return
	}

	if messaging.Message.Id != "" {
		if hook.handleInstagramMessage != nil {
			hook.handleInstagramMessage(ctx, messaging.Sender.Id, messaging.Recipient.Id, messaging.Timestamp, messaging.Message)
		}
		return
	}

	if messaging.Postback.Id != "" {
		if hook.handleInstagramPostback != nil {
			hook.handleInstagramPostback(ctx, messaging.Sender.Id, messaging.Recipient.Id, messaging.Timestamp, messaging.Postback)
		}
		return
	}

	if messaging.Referral.Type != "" {
		if hook.handleInstagramReferral != nil {
			hook.handleInstagramReferral(ctx, messaging.Sender.Id, messaging.Recipient.Id, messaging.Timestamp, messaging.Referral)
		}
		return
	}

	// TODO: seen and reaction, see https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook
}
