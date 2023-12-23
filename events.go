package gometawebhooks

import (
	"context"
	"fmt"
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
	Messaging []Messaging `json:"messaging,omitempty"`
	Changes   []Change    `json:"changes,omitempty"`
}

type Event struct {
	Object Object  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type EntryHandler interface {
	Entry(ctx context.Context, object Object, entry Entry)
}

func (h defaultHandler) Entry(ctx context.Context, object Object, entry Entry) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		h.hooks.changes(ctx, object, entry)
	}()

	go func() {
		defer wg.Done()

		h.hooks.messaging(ctx, object, entry)
	}()

	wg.Wait()
}
