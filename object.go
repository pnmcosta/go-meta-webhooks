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
	ErrObjectRequired     = errors.New("object required")
	ErrObjectNotSupported = errors.New("object not supported")

	supportedObjects = map[string]Object{
		"instagram": Instagram,
	}
)

func (t Object) String() string {
	return string(t)
}

func (t *Object) FromString(status string) Object {
	return supportedObjects[status]
}

func (t Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Object) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		return ErrObjectRequired
	}

	if _, ok := supportedObjects[s]; !ok {
		return fmt.Errorf("'%s': %w", s, ErrObjectNotSupported)
	}

	*t = t.FromString(s)
	return nil
}
