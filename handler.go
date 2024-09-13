package gometawebhooks

import (
	"context"
	"io"
	"net/http"
	"sync"
)

var _ EntryHandler = (*defaultHandler)(nil)
var _ ChangesHandler = (*defaultHandler)(nil)
var _ MessagingHandler = (*defaultHandler)(nil)

type defaultHandler struct {
	hooks *Webhooks
}

// Handles Meta Webhooks POST requests, verifies signature if secret is supplied, validates and parses Event payload.
func (hooks Webhooks) HandleRequest(ctx context.Context, r *http.Request) (Event, []byte, error) {
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

	if err := hooks.Verify(payload, headers); err != nil {
		return event, payload, err
	}

	if err := hooks.Validate(payload); err != nil {
		return event, payload, err
	}

	event, err = hooks.Parse(payload)
	if err != nil {
		return event, payload, err
	}

	var wg sync.WaitGroup
out:
	for _, entry := range event.Entry {
		select {
		case <-ctx.Done():
			break out
		default:
			wg.Add(1)

			entry := entry

			go func() {
				defer wg.Done()

				hooks.entryHandler.Entry(ctx, event.Object, entry)
			}()
		}
	}

	wg.Wait()

	select {
	case <-ctx.Done():
		return event, payload, ctx.Err()
	default:
		return event, payload, nil
	}
}
