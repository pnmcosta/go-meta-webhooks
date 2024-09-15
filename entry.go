package gometawebhooks

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidHTTPMethod         = fmt.Errorf("invalid HTTP Method: %w", ErrWebhooks)
	ErrReadBodyPayload           = fmt.Errorf("error reading body payload: %w", ErrWebhooks)
	ErrMissingHubSignatureHeader = fmt.Errorf("missing signature value: %w", ErrWebhooks)
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

type EntryHandler interface {
	Entry(ctx context.Context, object Object, entry Entry) error
}

func (h Webhooks) Entry(ctx context.Context, object Object, entry Entry) error {
	g := new(errgroup.Group)
	g.SetLimit(2)

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return h.changes(ctx, object, entry)
		}
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return h.messaging(ctx, object, entry)
		}
	})

	return g.Wait()
}
