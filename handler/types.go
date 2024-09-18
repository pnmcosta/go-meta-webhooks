package handler

import gometawebhooks "github.com/pnmcosta/go-meta-webhooks"

type (
	Change            = gometawebhooks.Change
	Event             = gometawebhooks.Event
	Messaging         = gometawebhooks.Messaging
	Message           = gometawebhooks.Message
	Entry             = gometawebhooks.Entry
	Option            = gometawebhooks.Option
	Object            = gometawebhooks.Object
	Postback          = gometawebhooks.Postback
	Referral          = gometawebhooks.Referral
	Attachment        = gometawebhooks.Attachment
	AttachmentPayload = gometawebhooks.AttachmentPayload

	MessagingHeader               = gometawebhooks.MessagingHeader
	MessagingMessage              = gometawebhooks.MessagingMessage
	MessagingPostback             = gometawebhooks.MessagingPostback
	MessagingReferral             = gometawebhooks.MessagingReferral
	InstagramHandler              = gometawebhooks.InstagramHandler
	InstagramChangesHandler       = gometawebhooks.InstagramChangesHandler
	InstagramMentionHandler       = gometawebhooks.InstagramMentionHandler
	InstagramMessageHandler       = gometawebhooks.InstagramMessageHandler
	InstagramMessagingHandler     = gometawebhooks.InstagramMessagingHandler
	InstagramPostbackHandler      = gometawebhooks.InstagramPostbackHandler
	InstagramReferralHandler      = gometawebhooks.InstagramReferralHandler
	InstagramStoryInsightsHandler = gometawebhooks.InstagramStoryInsightsHandler

	Mention       = gometawebhooks.Mention
	StoryInsights = gometawebhooks.StoryInsights
)

var Options = gometawebhooks.Options

const (
	Instagram = gometawebhooks.Instagram
)
