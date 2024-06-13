package gometawebhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
)

var _ EntryHandler = (*defaultHandler)(nil)
var _ ChangesHandler = (*defaultHandler)(nil)
var _ MessagingHandler = (*defaultHandler)(nil)

type defaultHandler struct {
	hooks *Webhooks
}

// Handles Meta Webhooks POST requests, verifies signature if secret is supplied, validates and parses Event payload.
func (hooks Webhooks) Handle(ctx context.Context, r *http.Request) (Event, error) {
	defer func() {
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
	}()

	var event Event

	if r.Method != http.MethodPost {
		return event, ErrInvalidHTTPMethod
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return event, wrapErr(err, ErrReadBodyPayload)
	}

	// normalize header keys
	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[strings.ToLower(strings.ReplaceAll(k, "-", "_"))] = v[0]
		}
	}

	// If we have a Secret set, we should check the MAC
	// https://developers.facebook.com/docs/messenger-platform/webhooks#validate-payloads
	if len(hooks.secret) > 0 {
		signature := headers["x_hub_signature_256"]
		if len(signature) == 0 {
			return event, ErrMissingHubSignatureHeader
		}
		mac := hmac.New(sha256.New, []byte(hooks.secret))
		mac.Write(payload)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		if len(signature) <= 8 || !hmac.Equal([]byte(signature[7:]), []byte(expectedMAC)) {
			return event, ErrHMACVerificationFailed
		}
	}

	var pl interface{}
	if err := json.Unmarshal(payload, &pl); err != nil {
		return event, wrapErr(err, ErrParsingPayload)
	}

	if err := hooks.validate(pl); err != nil {
		return event, err
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return event, wrapErr(err, ErrParsingEvent)
	}

	var wg sync.WaitGroup
out:
	for _, entry := range event.Entry {
		select {
		case <-ctx.Done():
			break out
		default:
		}

		wg.Add(1)

		entry := entry

		go func() {
			defer wg.Done()

			hooks.entryHandler.Entry(ctx, event.Object, entry)
		}()
	}

	wg.Wait()

	select {
	case <-ctx.Done():
		return event, ctx.Err()
	default:
	}

	return event, nil
}
