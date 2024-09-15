package gometawebhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const (
	HeaderSignatureName = "X-Hub-Signature-256"
)

var (
	ErrMissingHubSignatureHeader = fmt.Errorf("missing signature value: %w", ErrWebhooks)
	ErrHMACVerificationFailed    = fmt.Errorf("HMAC verification failed: %w", ErrWebhooks)
	ErrParsingPayload            = fmt.Errorf("error parsing payload: %w", ErrWebhooks)
	ErrParsingEvent              = fmt.Errorf("error parsing event: %w", ErrWebhooks)
)

func (hooks Webhooks) ParsePayload(body []byte) (Event, error) {
	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		return event, wrapErr(err, ErrParsingEvent)
	}
	return event, nil
}

func (hooks Webhooks) ValidatePayload(body []byte) error {
	if validationSchema == nil {
		return ErrMissingSchema
	}

	var pl interface{}
	if err := json.Unmarshal(body, &pl); err != nil {
		return wrapErr(err, ErrParsingPayload)
	}

	if err := validationSchema.Validate(pl); err != nil {
		return wrapErr(err, ErrInvalidPayload)
	}

	return nil
}

func (hooks Webhooks) VerifyPayload(body []byte, headers map[string]string) error {
	// If we have a Secret set, we should check the MAC
	// https://developers.facebook.com/docs/messenger-platform/webhooks#validate-payloads
	if len(hooks.secret) == 0 {
		return nil
	}

	signature := headers[hooks.headerSigName]
	if len(signature) == 0 {
		return fmt.Errorf("missing %s Header: %w", hooks.headerSigName, ErrMissingHubSignatureHeader)
	}

	mac := hmac.New(sha256.New, []byte(hooks.secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	if len(signature) <= 8 || !hmac.Equal([]byte(signature[7:]), []byte(expectedMAC)) {
		return ErrHMACVerificationFailed
	}
	return nil
}