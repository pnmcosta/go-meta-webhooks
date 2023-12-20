package gometawebhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	ErrInvalidHTTPMethod         = fmt.Errorf("invalid HTTP Method: %w", ErrWebhooks)
	ErrReadBodyPayload           = fmt.Errorf("error reading body payload: %w", ErrWebhooks)
	ErrMissingHubSignatureHeader = fmt.Errorf("missing x-hub-signature-256 Header: %w", ErrWebhooks)
	ErrHMACVerificationFailed    = fmt.Errorf("HMAC verification failed: %w", ErrWebhooks)
	ErrParsingPayload            = fmt.Errorf("error parsing payload: %w", ErrWebhooks)
	ErrParsingEvent              = fmt.Errorf("error parsing event: %w", ErrWebhooks)
)

type Entry struct {
	Id        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
	Changes   []Change    `json:"changes"`
}

type Event struct {
	Object Object  `json:"object"`
	Entry  []Entry `json:"entry"`
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

	fns := []func(context.Context, Object, Entry){
		hooks.changes,
		hooks.messaging,
	}

	for _, entry := range event.Entry {
		entry := entry

		var wg sync.WaitGroup
		wg.Add(len(fns))

		for _, fn := range fns {
			fn := fn

			go func(entry Entry) {
				defer wg.Done()

				fn(ctx, event.Object, entry)
			}(entry)
		}

		wg.Wait()
	}

	return event, nil
}
