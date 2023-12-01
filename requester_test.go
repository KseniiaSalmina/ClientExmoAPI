package main

import (
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := &http.Client{}
	expected := &Client{client: client}
	result := NewClient(client)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result: got %v, got %v", *result, *expected)
	}
}

func TestClient_GetRequest(t *testing.T) {
	exmo := NewExmo()
	type testData struct {
		url         string
		method      string
		body        io.Reader
		expectedErr bool
	}
	start := strconv.Itoa(int(time.Now().Add(-time.Hour * 24).Unix()))
	end := strconv.Itoa(int(time.Now().Unix()))

	testCases := []testData{
		{url: "not url", method: "POST", body: io.Reader(nil), expectedErr: true},
		{url: "https://api.exmo.com/v1.1/ticker", method: "POST", body: io.Reader(nil), expectedErr: false},
		{url: "https://api.exmo.com/v1.1/currency", method: "POST", body: io.Reader(nil), expectedErr: false},
		{url: "https://api.exmo.com/v1.1/trades", method: "POST", body: strings.NewReader(`pair=BTC_USD`), expectedErr: false},
		{url: "https://api.exmo.com/v1.1/candles_history?symbol=BTC_USD&resolution=30&from=" + start + "&to=" + end, method: "GET", body: io.Reader(nil), expectedErr: false},
	}

	for _, tc := range testCases {
		result, err := exmo.requester.GetRequest(tc.method, tc.url, tc.body)
		if tc.expectedErr {
			if err == nil {
				t.Errorf("url: %v: expected error, got nil", tc.url)
			}
			if result != nil {
				t.Errorf("url: %v: unexpected result %v, want nil", tc.url, result)
			}
		}
		if !tc.expectedErr {
			if err != nil {
				t.Errorf("url: %v: unexpected error %v", tc.url, err)
			}
			if result == nil {
				t.Errorf("url: %v: unexpected nil result", tc.url)
			}
		}
	}
}
