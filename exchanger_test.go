package main

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewExmo(t *testing.T) {
	expected := &Exmo{client: &http.Client{}, url: "https://api.exmo.com/v1.1", isTest: false, requester: NewClient(&http.Client{})}
	result := NewExmo()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result: got %v, got %v", *result, *expected)
	}
}

func TestWithClient(t *testing.T) {
	client := &http.Client{}
	expected := &Exmo{client: client, url: "https://api.exmo.com/v1.1", isTest: false, requester: NewClient(client)}
	result := NewExmo(WithClient(client))

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result: got %v, got %v", *result, *expected)
	}
}

func TestWithURL(t *testing.T) {
	url := "https://www.test.com"
	expected := &Exmo{client: &http.Client{}, url: url, isTest: false, requester: NewClient(&http.Client{})}
	result := NewExmo(WithURL(url))
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result: got %v, got %v", *result, *expected)
	}
}

func TestTest(t *testing.T) {
	expected := &Exmo{client: &http.Client{}, url: "https://api.exmo.com/v1.1", isTest: true, requester: &MockClient{}}
	result := NewExmo(Test())
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected result: got %v, got %v", *result, *expected)
	}
}

func TestExmo_GetTicker(t *testing.T) {
	exmo := NewExmo(Test())
	expectedPairs := []string{"ADA_BTC", "ADA_USD"}
	result, err := exmo.GetTicker()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected result, got nil")
	} else {
		for _, pair := range expectedPairs {
			if _, ok := result[pair]; !ok {
				t.Errorf("unexpected result: got result[%v] false, want true", pair)
			}
		}
	}
}

func TestExmo_GetTrades(t *testing.T) {
	exmo := NewExmo(Test())
	expectedPairs := []string{"ADA_BTC", "ADA_USD"}

	for _, pair := range expectedPairs {
		result, err := exmo.GetTrades(pair)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Errorf("expected result, got nil")
		}
	}

	result2, err2 := exmo.GetTrades(expectedPairs...)
	if err2 != nil {
		t.Errorf("unexpected error: %v", err2)
	}
	if result2 == nil {
		t.Errorf("expected result, got nil")
	} else {
		for _, pair := range expectedPairs {
			if _, ok := result2[pair]; !ok {
				t.Errorf("unexpected result: got result[%v] false, want true", pair)
			}
		}
	}
}

func TestExmo_GetOrderBook(t *testing.T) {
	exmo := NewExmo(Test())
	expectedPairs := []string{"ADA_BTC", "ADA_USD"}

	for _, pair := range expectedPairs {
		result, err := exmo.GetOrderBook(30, pair)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Errorf("expected result, got nil")
		}
	}

	result2, err2 := exmo.GetOrderBook(30, expectedPairs...)
	if err2 != nil {
		t.Errorf("unexpected error: %v", err2)
	}
	if result2 == nil {
		t.Errorf("expected result, got nil")
	} else {
		for _, pair := range expectedPairs {
			if _, ok := result2[pair]; !ok {
				t.Errorf("unexpected result: got result[%v] false, want true", pair)
			}
		}
	}
}

func TestExmo_GetCurrencies(t *testing.T) {
	exmo := NewExmo(Test())
	expectedCurrencies := Currencies{"ADA_BTC", "ADA_USD"}

	result, err := exmo.GetCurrencies()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expectedCurrencies) {
		t.Errorf("unexpected result, got %v, want %v", result, expectedCurrencies)
	}
}

func TestExmo_GetCandlesHistory(t *testing.T) {
	exmo := NewExmo(Test())
	pair := "ADA_BTC"

	result, err := exmo.GetCandlesHistory(pair, 30, time.Unix(1701367794, 0), time.Unix(1701367795, 0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if reflect.DeepEqual(result, CandlesHistory{}) {
		t.Errorf("expected result %v, got nil, ", result)
	}
}

func TestExmo_GetClosePrice(t *testing.T) {
	exmo := NewExmo(Test())
	pair := "ADA_BTC"

	result, err := exmo.GetClosePrice(pair, 30, time.Unix(1701367794, 0), time.Unix(1701367795, 0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected result %v, got nil, ", result)
	}
}
