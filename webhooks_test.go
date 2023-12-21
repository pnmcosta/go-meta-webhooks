package gometawebhooks_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
)

type hookScenario struct {
	name             string
	method           string
	url              string
	headers          map[string]string
	options          func(scenario *hookScenario) []gometawebhooks.Option
	body             io.Reader
	expected         interface{}
	expectErr        error
	expectedHandlers map[string]int
	handled          []string
	timeout          time.Duration
}

func (scenario *hookScenario) test(t *testing.T, f func(t *testing.T)) {
	var name = scenario.name
	if name == "" {
		name = fmt.Sprintf("%s:%s", scenario.method, scenario.url)
	}

	t.Run(name, f)
}

func (scenario *hookScenario) setup(t *testing.T) (*gometawebhooks.Webhooks, *http.Request) {
	if scenario.timeout == 0 {
		scenario.timeout = 100 * time.Millisecond
	}

	var options []gometawebhooks.Option
	if scenario.options != nil {
		options = scenario.options(scenario)
	}

	hooks, err := gometawebhooks.NewWebhooks(options...)
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

	if scenario.expectedHandlers != nil {
		counter := map[string]int{}
		for _, s := range scenario.handled {
			counter[s]++
		}

		if !maps.Equal(counter, scenario.expectedHandlers) {
			t.Errorf("Expected handlers %v, but got %v", scenario.expectedHandlers, counter)
		}
	}

	if !reflect.DeepEqual(result, scenario.expected) {
		t.Errorf("Expected %v, but got %v", scenario.expected, result)
	}
}

func (scenario *hookScenario) trigger(event string) {
	scenario.handled = append(scenario.handled, event)
}

func genHmac(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
