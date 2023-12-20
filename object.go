package gometawebhooks

import "encoding/json"

type Object string

const (
	Instagram Object = "instagram"
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
		return err
	}
	*t = t.FromString(s)
	return nil
}
