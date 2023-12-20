package gometawebhooks

import "context"

func (hook Webhooks) defaultOnInstagramMessaging(ctx context.Context, entry Entry, messaging Messaging) {
	if messaging.Message.IsEcho {
		return
	}

	if messaging.Message.Id != "" {
		if hook.onInstagramMessage != nil {
			hook.onInstagramMessage(ctx, messaging.Sender.Id, messaging.Recipient.Id, messaging.Timestamp, messaging.Message)
		}
		return
	}

	if messaging.Postback.Id != "" {
		if hook.onInstagramPostback != nil {
			hook.onInstagramPostback(ctx, messaging.Sender.Id, messaging.Recipient.Id, messaging.Timestamp, messaging.Postback)
		}
		return
	}

	if messaging.Referral.Type != "" {
		if hook.onInstagramReferral != nil {
			hook.onInstagramReferral(ctx, messaging.Sender.Id, messaging.Recipient.Id, messaging.Timestamp, messaging.Referral)
		}
		return
	}

	// TODO: seen and reaction, see https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook
}
