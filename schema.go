package gometawebhooks

import (
	_ "embed"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var (
	//go:embed schema.json
	embedSchema      string
	validationSchema *jsonschema.Schema

	ErrSchemaCompile   = fmt.Errorf("failed to compile schema: %w", ErrWebhooks)
	ErrSchemaNotSet    = fmt.Errorf("missing embedded schema: %w", ErrWebhooks)
	ErrPayloadValidate = fmt.Errorf("failed to validate payload: %w", ErrWebhooks)
)

func (hook *Webhooks) compileSchema() error {
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

func (hook Webhooks) validate(payload interface{}) error {
	if validationSchema == nil {
		return ErrSchemaNotSet
	}

	if err := validationSchema.Validate(payload); err != nil {
		return wrapErr(err, ErrPayloadValidate)
	}

	return nil
}
