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
	g := new(errgroup.Group)
	g.SetLimit(1)
	for _, entry := range event.Entry {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return h.Entry(ctx, event.Object, entry)
			}
		})
	}
	return g.Wait()
}
