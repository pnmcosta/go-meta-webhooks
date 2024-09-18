package gometawebhooks

import (
	_ "embed"
	"errors"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var (
	//go:embed schema.json
	embedSchema      string
	validationSchema *jsonschema.Schema

	ErrSchemaCompile = errors.New("failed to compile schema")
	ErrMissingSchema = errors.New("missing embedded schema")

	mu sync.RWMutex
)

func (hooks *Webhooks) compileSchema() error {
	mu.Lock()
	defer mu.Unlock()
	if validationSchema != nil {
		return nil
	}

	schema, err := jsonschema.CompileString("schema.json", embedSchema)
	if err != nil {
		return wrapErr(err, ErrSchemaCompile)
	}
	validationSchema = schema
	return nil
}
