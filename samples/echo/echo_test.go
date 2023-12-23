package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestVerify(t *testing.T) {
	e, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}

	expectedResponse := "challenge_response"

	q := make(url.Values)
	q.Set("hub.mode", "subscribe")
	q.Set("hub.challenge", expectedResponse)
	q.Set("hub.verify_token", MetaWebhooksToken)

	req := httptest.NewRequest(http.MethodGet, MetaWebhooksRoute+"?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	body := rec.Body.String()
	if !strings.Contains(body, expectedResponse) {
		t.Errorf("Cannot find %v in response body \n%v", expectedResponse, body)
	}
}

func TestInstagramMessage(t *testing.T) {
	e, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}

	body := strings.NewReader(`{
		"object": "instagram",
		"entry": [
		  {
			"id": "123",
			"time": 1569262486134,
			"messaging": [
			  {
				"sender": {
				  "id": "567"
				},
				"recipient": {
				  "id": "123"
				},
				"timestamp": 1569262485349,
				"message": {
				  "mid": "890",
				  "text": "Text in message"
				}
			  }
			]
		  }
		]
	  }`)
	req := httptest.NewRequest(http.MethodPost, MetaWebhooksRoute, body)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}
