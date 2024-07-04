package gometawebhooks_test

import (
	"bytes"
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
	"sync"
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
	bodyBytes        []byte
	expected         interface{}
	expectErr        error
	expectedHandlers map[string]int
	timeout          time.Duration

	handled []string
	mutex   *sync.RWMutex
}

func (scenario *hookScenario) test(t *testing.T, f func(t *testing.T)) {
	var name = scenario.name
	if name == "" {
		name = fmt.Sprintf("%s:%s", scenario.method, scenario.url)
	}

	t.Run(name, f)
}

func (scenario *hookScenario) setup(t *testing.T) (*gometawebhooks.Webhooks, *http.Request) {
	scenario.mutex = &sync.RWMutex{}

	if scenario.timeout == 0 {
		scenario.timeout = 7 * time.Millisecond
	}

	var options []gometawebhooks.Option
	if scenario.options != nil {
		options = scenario.options(scenario)
	}

	hooks, err := gometawebhooks.New(options...)
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

	// Read and store the body for later use
	scenario.bodyBytes, err = io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("failed to read request body: %v", err)
	}
	req.Body.Close() // close the original body

	// Restore the body to the request for the handler to read
	req.Body = io.NopCloser(bytes.NewBuffer(scenario.bodyBytes))

	return hooks, req
}

func (scenario *hookScenario) assert(t *testing.T, result interface{}, payload []byte, err error) {
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

	if scenario.body != nil {
		if !bytes.Equal(scenario.bodyBytes, payload) {
			t.Errorf("Expected body %v, but got %v", string(scenario.bodyBytes), string(payload))
		}
	}
}

func (scenario *hookScenario) trigger(event string) {
	scenario.mutex.Lock()
	defer scenario.mutex.Unlock()
	scenario.handled = append(scenario.handled, event)
}

func genHmac(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
