package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

var (
	ErrReadBodyPayload   = errors.New("error reading body payload")
	ErrInvalidHTTPMethod = errors.New("invalid HTTP Method")
)

type DefaultHandler interface {
	gometawebhooks.WebhooksHandler

	HandleRequest(ctx context.Context, r *http.Request) (Event, []byte, error)
	HandleVerify(r *http.Request) (string, error)
}

var _ DefaultHandler = (*defaultHandler)(nil)

type defaultHandler struct {
	*gometawebhooks.Webhooks
}

func New(opts ...Option) (*defaultHandler, error) {
	if len(opts) == 0 {
		opts = append(opts, Options.CompileSchema())
	}

	hooks, err := gometawebhooks.New(opts...)
	if err != nil {
		return nil, err
	}

	return &defaultHandler{hooks}, nil
}

// Handles Meta Webhooks POST requests, verifies signature if secret is supplied, validates and parses Event payload.
func (hooks defaultHandler) HandleRequest(ctx context.Context, r *http.Request) (Event, []byte, error) {
	defer func() {
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
	}()

	var event Event

	if r.Method != http.MethodPost {
		return event, []byte{}, ErrInvalidHTTPMethod
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return event, payload, wrapErr(err, ErrReadBodyPayload)
	}

	// normalize header keys
	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	if err := hooks.VerifyPayload(payload, headers); err != nil {
		return event, payload, err
	}

	if err := hooks.ValidatePayload(payload); err != nil {
		return event, payload, err
	}

	event, err = hooks.ParsePayload(payload)
	if err != nil {
		return event, payload, err
	}

	err = hooks.Handle(ctx, event)
	return event, payload, err
}

// Verify Meta Webhooks GET requests, when subscribing on App dashboard to objects and fields.
func (hooks defaultHandler) HandleVerify(r *http.Request) (string, error) {
	if r.Method != http.MethodGet {
		return "", ErrInvalidHTTPMethod
	}

	q := r.URL.Query()
	return hooks.VerifyToken(map[string]string{
		"hub.mode":         q.Get("hub.mode"),
		"hub.verify_token": q.Get("hub.verify_token"),
		"hub.challenge":    q.Get("hub.challenge"),
	})
}

func wrapErr(err, target error) error {
	return fmt.Errorf("%s: %w", err, target)
}
