package gometawebhooks

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Object string

const (
	Instagram Object = "instagram"
)

var (
	ErrObjectNotSupported = fmt.Errorf("object not supported: %w", ErrWebhooks)
)

func (t Object) String() string {
	return string(t)
}

func (t *Object) FromString(status string) Object {
	return map[string]Object{
		"instagram": Instagram,
	}[status]
}

func (t Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Object) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return wrapErr(err, ErrObjectNotSupported)
	}
	*t = t.FromString(s)
	if *t == "" {
		return wrapErr(errors.New(s), ErrObjectNotSupported)
	}

	return nil
}
