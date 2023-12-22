package gometawebhooks

import (
	"encoding/json"
	"fmt"
)

type Object string

const (
	Instagram Object = "instagram"
)

var (
	ErrObjectRequired     = fmt.Errorf("object required: %w", ErrWebhooks)
	ErrObjectNotSupported = fmt.Errorf("object not supported: %w", ErrWebhooks)
)

func (t Object) String() string {
	return string(t)
}

func (t *Object) FromString(status string) Object {
	return Object(status)
}

func (t Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Object) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == "" {
		return ErrObjectRequired
	}

	*t = t.FromString(s)
	return nil
}
