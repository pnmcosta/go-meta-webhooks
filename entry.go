package gometawebhooks

import (
	"context"

	"golang.org/x/sync/errgroup"
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
	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(2)

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			return h.changes(ctx, object, entry)
		}
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			return h.messaging(ctx, object, entry)
		}
	})

	return g.Wait()
}
