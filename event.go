package gometawebhooks

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Event struct {
	Object Object  `json:"object"`
	Entry  []Entry `json:"entry"`
}

func (h Webhooks) Handle(ctx context.Context, event Event) error {
	if len(event.Entry) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(len(event.Entry))
	for _, entry := range event.Entry {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return context.Cause(ctx)
			default:
				return h.entry(ctx, event.Object, entry)
			}
		})
	}
	return g.Wait()
}
