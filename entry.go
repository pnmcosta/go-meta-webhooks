package gometawebhooks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var (
	ErrParsingEntry = errors.New("parsing entry")
)

type Entry struct {
	Id        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging,omitempty"`
	Changes   []Change    `json:"changes,omitempty"`
}

func (t *Entry) UnmarshalJSON(b []byte) error {
	type Alias Entry
	var entry Alias
	if err := json.Unmarshal(b, &entry); err != nil {
		return err
	}

	if entry.Id == "" {
		return fmt.Errorf("missing 'id' field: %w", ErrParsingEntry)
	}

	if entry.Time == 0 {
		return fmt.Errorf("missing 'time' field: %w", ErrParsingEntry)
	}

	*t = Entry(entry)
	return nil
}

type EntryHandler interface {
	Entry(context.Context, Object, Entry) error
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
