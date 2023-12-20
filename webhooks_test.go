package gometawebhooks_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

type hookScenario struct {
	name      string
	method    string
	url       string
	headers   map[string]string
	options   []gometawebhooks.Option
	body      io.Reader
	expected  interface{}
	expectErr error
}

func (scenario *hookScenario) test(t *testing.T, f func(t *testing.T)) {
	var name = scenario.name
	if name == "" {
		name = fmt.Sprintf("%s:%s", scenario.method, scenario.url)
	}

	t.Run(name, f)
}

func (scenario *hookScenario) init(t *testing.T) (*gometawebhooks.Webhooks, *http.Request) {
	hooks, err := gometawebhooks.NewWebhooks(scenario.options...)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	method := scenario.method
	if method == "" {
		method = http.MethodGet
	}

	url := scenario.url
	if url == "" {
		url = "/webhooks/meta"
	}

	req := httptest.NewRequest(method, url, scenario.body)
	for k, v := range scenario.headers {
		req.Header.Set(k, v)
	}

	return hooks, req
}

func (scenario *hookScenario) assert(t *testing.T, result interface{}, err error) {
	if scenario.expectErr != nil {
		if err == nil {
			t.Errorf("Expected an error, but got none.")
		}

		if !errors.Is(err, scenario.expectErr) {
			t.Errorf("Expected error %v, but got %v.", scenario.expectErr, err)
		}
		return
	}

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if !reflect.DeepEqual(result, scenario.expected) {
		t.Errorf("Expected %v, but got %v", scenario.expected, result)
	}
}

func genHmac(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
