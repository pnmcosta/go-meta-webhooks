package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"log"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

const (
	MetaWebhooksToken = "my-webhook-token"
	MetaWebhooksRoute = "/webhooks/meta"
)

type handler struct {
}

var _ gometawebhooks.InstagramHandler = (*handler)(nil)

func main() {
	mux, err := setup()
	if err != nil {
		log.Print(err)
		return
	}

	if err := http.ListenAndServe("127.0.0.1:1323", mux); err != nil {
		log.Print(err)
	}
}

func setup() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	handler := handler{}

	hooks, err := gometawebhooks.New(
		// gometawebhooks.Options.Secret("my-app-secret"),
		gometawebhooks.Options.Token(MetaWebhooksToken),
		gometawebhooks.Options.InstagramHandler(handler),
	)
	if err != nil {
		return mux, err
	}

	mux.HandleFunc(MetaWebhooksRoute, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			challenge, err := hooks.Verify(r)
			if err != nil {
				log.Print(err)
				return
			}

			w.WriteHeader(http.StatusOK)
			io.WriteString(w, challenge)
		case http.MethodPost:
			_, err := hooks.Handle(r.Context(), r)
			if err != nil {
				log.Print(err)
				return
			}

			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "unsupported")
		}
	})

	return mux, nil
}

func (h handler) InstagramMessage(ctx context.Context, sender, recipient string, sent time.Time, message gometawebhooks.Message) {
	log.Printf("instagram message: %v, %v, %v, %v", sender, recipient, sent, message)
}

func (h handler) InstagramPostback(ctx context.Context, sender string, recipient string, sent time.Time, postback gometawebhooks.Postback) {
	log.Printf("instagram postback: %v, %v, %v, %v", sender, recipient, sent, postback)
}

func (h handler) InstagramReferral(ctx context.Context, sender string, recipient string, sent time.Time, referral gometawebhooks.Referral) {
	log.Printf("instagram referral: %v, %v, %v, %v", sender, recipient, sent, referral)
}

func (h handler) InstagramStoryInsights(ctx context.Context, entryId string, entryTime time.Time, storyInsights gometawebhooks.StoryInsightsFieldValue) {
	log.Printf("instagram story insights: %v, %v, %v", entryId, entryTime, storyInsights)
}

func (h handler) InstagramMention(ctx context.Context, entryId string, entryTime time.Time, mention gometawebhooks.MentionsFieldValue) {
	log.Printf("instagram mention: %v, %v, %v", entryId, entryTime, mention)
}
