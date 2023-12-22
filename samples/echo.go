package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

const (
	MetaWebhookToken = "my-webhook-token"
)

type Handler struct {
	logger echo.Logger
}

func main() {
	e := echo.New()

	handler := Handler{
		logger: e.Logger,
	}

	hooks, err := gometawebhooks.NewWebhooks(
		// gometawebhooks.Options.Secret("my-app-secret"),
		gometawebhooks.Options.Token(MetaWebhookToken),
		gometawebhooks.Options.InstagramMessageHandler(handler),
	)

	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	e.GET("/webhooks/meta", func(c echo.Context) error {
		challenge, err := hooks.Verify(c.Request())
		if err != nil {
			e.Logger.Error(err)
			return err
		}

		return c.String(http.StatusOK, challenge)
	})
	e.POST("/webhooks/meta", func(c echo.Context) error {
		_, err := hooks.Handle(c.Request().Context(), c.Request())
		if err != nil {
			e.Logger.Error(err)
			return err
		}

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

func (h Handler) InstagramMessage(ctx context.Context, sender, recipient string, time time.Time, message gometawebhooks.Message) {
	h.logger.Infof("instagram message from %s to %s at %v with payload: %v", sender, recipient, time, message)
}
