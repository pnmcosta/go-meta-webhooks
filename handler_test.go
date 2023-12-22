package gometawebhooks_test

import (
	"context"
	"time"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

type testHandler struct {
	run func()
}

// Entry implements gometawebhooks.EntryHandler.
func (h testHandler) Entry(ctx context.Context, object gometawebhooks.Object, entry gometawebhooks.Entry) {
	h.run()
}

// Changes implements gometawebhooks.ChangesHandler.
func (h testHandler) Changes(ctx context.Context, object gometawebhooks.Object, entry gometawebhooks.Entry, change gometawebhooks.Change) {
	h.run()
}

// Messaging implements gometawebhooks.MessagingHandler.
func (h testHandler) Messaging(ctx context.Context, object gometawebhooks.Object, entryId string, entryTime time.Time, messaging gometawebhooks.Messaging) {
	h.run()
}

// InstagramMention implements gometawebhooks.InstagramMentionHandler.
func (h testHandler) InstagramMention(ctx context.Context, entryId string, entryTime time.Time, mention gometawebhooks.MentionsFieldValue) {
	h.run()
}

// InstagramStoryInsights implements gometawebhooks.InstagramStoryInsightsHandler.
func (h testHandler) InstagramStoryInsights(ctx context.Context, entryId string, entryTime time.Time, storyInsights gometawebhooks.StoryInsightsFieldValue) {
	h.run()
}

// InstagramMessage implements gometawebhooks.InstagramMessageHandler.
func (h testHandler) InstagramMessage(ctx context.Context, sender string, recipient string, sent time.Time, message gometawebhooks.Message) {
	h.run()
}

// InstagramPostback implements gometawebhooks.InstagramPostbackHandler.
func (h testHandler) InstagramPostback(ctx context.Context, sender string, recipient string, sent time.Time, postback gometawebhooks.Postback) {
	h.run()
}

// InstagramReferral implements gometawebhooks.InstagramReferralHandler.
func (h testHandler) InstagramReferral(ctx context.Context, sender string, recipient string, sent time.Time, referral gometawebhooks.Referral) {
	h.run()
}

var _ gometawebhooks.EntryHandler = (*testHandler)(nil)
var _ gometawebhooks.ChangesHandler = (*testHandler)(nil)
var _ gometawebhooks.MessagingHandler = (*testHandler)(nil)
var _ gometawebhooks.InstagramMentionHandler = (*testHandler)(nil)
var _ gometawebhooks.InstagramStoryInsightsHandler = (*testHandler)(nil)
var _ gometawebhooks.InstagramMessageHandler = (*testHandler)(nil)
var _ gometawebhooks.InstagramPostbackHandler = (*testHandler)(nil)
var _ gometawebhooks.InstagramReferralHandler = (*testHandler)(nil)
