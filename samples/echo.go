package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

const (
	MetaWebhooksToken = "my-webhook-token"
	MetaWebhooksRoute = "/webhooks/meta"
)

type Handler struct {
	logger echo.Logger
}

func main() {
	e, err := setup()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

func setup() (*echo.Echo, error) {
	e := echo.New()

	handler := Handler{
		logger: e.Logger,
	}

	hooks, err := gometawebhooks.NewWebhooks(
		// gometawebhooks.Options.Secret("my-app-secret"),
		gometawebhooks.Options.Token(MetaWebhooksToken),
		gometawebhooks.Options.InstagramHandler(handler),
	)
	if err != nil {
		return e, err
	}

	e.GET(MetaWebhooksRoute, func(c echo.Context) error {
		challenge, err := hooks.Verify(c.Request())
		if err != nil {
			e.Logger.Error(err)
			return err
		}

		return c.String(http.StatusOK, challenge)
	})

	e.POST(MetaWebhooksRoute, func(c echo.Context) error {
		_, err := hooks.Handle(c.Request().Context(), c.Request())
		if err != nil {
			e.Logger.Error(err)
			return err
		}

		return c.NoContent(http.StatusOK)
	})

	return e, nil
}

func (h Handler) InstagramMessage(ctx context.Context, sender, recipient string, sent time.Time, message gometawebhooks.Message) {
	h.logger.Infof("instagram message: %v, %v, %v, %v", sender, recipient, sent, message)
}

func (h Handler) InstagramPostback(ctx context.Context, sender string, recipient string, sent time.Time, postback gometawebhooks.Postback) {
	h.logger.Infof("instagram postback: %v, %v, %v, %v", sender, recipient, sent, postback)
}

func (h Handler) InstagramReferral(ctx context.Context, sender string, recipient string, sent time.Time, referral gometawebhooks.Referral) {
	h.logger.Infof("instagram referral: %v, %v, %v, %v", sender, recipient, sent, referral)
}

func (h Handler) InstagramStoryInsights(ctx context.Context, entryId string, entryTime time.Time, storyInsights gometawebhooks.StoryInsightsFieldValue) {
	h.logger.Infof("instagram story insights: %v, %v, %v", entryId, entryTime, storyInsights)
}

func (h Handler) InstagramMention(ctx context.Context, entryId string, entryTime time.Time, mention gometawebhooks.MentionsFieldValue) {
	h.logger.Infof("instagram mention: %v, %v, %v", entryId, entryTime, mention)
}
