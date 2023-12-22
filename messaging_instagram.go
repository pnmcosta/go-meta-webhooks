package gometawebhooks

import (
	"context"
	"time"
)

// TODO: seen and reaction, see https://developers.facebook.com/docs/messenger-platform/instagram/features/webhook

type InstagramMessageHandler interface {
	InstagramMessage(ctx context.Context, sender, recipient string, sent time.Time, message Message)
}

type InstagramPostbackHandler interface {
	InstagramPostback(ctx context.Context, sender, recipient string, sent time.Time, postback Postback)
}

type InstagramReferralHandler interface {
	InstagramReferral(ctx context.Context, sender, recipient string, sent time.Time, referral Referral)
}

type InstagramMessagingHandler interface {
	InstagramMessageHandler
	InstagramPostbackHandler
	InstagramReferralHandler
}
