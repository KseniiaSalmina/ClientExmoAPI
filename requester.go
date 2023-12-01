package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{client: client}
}

func (c *Client) GetRequest(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("Client_GetRequest -> %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Client_GetRequest -> %w", err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Client_GetRequest -> %w", err)
	}

	return bodyText, nil
}

type MockClient struct {
}

func (m *MockClient) GetRequest(method string, url string, body io.Reader) ([]byte, error) {
	switch url {
	case "https://api.exmo.com/v1.1/ticker":
		return json.Marshal(Ticker{"ADA_BTC": TickerValue{}, "ADA_USD": TickerValue{}})

	case "https://api.exmo.com/v1.1/trades":
		if !reflect.DeepEqual(body, strings.NewReader("pair=ADA_BTC")) {
			return json.Marshal(Trades{"ADA_BTC": []Pair{}})
		}
		if !reflect.DeepEqual(body, strings.NewReader("pair=ADA_USD")) {
			return json.Marshal(Trades{"ADA_USD": []Pair{}})
		}
		return []byte("invalid test data"), errors.New("invalid test data")

	case "https://api.exmo.com/v1.1/order_book":
		if !reflect.DeepEqual(body, strings.NewReader("pair=ADA_BTC&limit=30")) {
			return json.Marshal(OrderBook{"ADA_BTC": OrderBookPair{}})
		}
		if !reflect.DeepEqual(body, strings.NewReader("pair=ADA_USD&limit=30")) {
			return json.Marshal(OrderBook{"ADA_USD": OrderBookPair{}})
		}
		return []byte("invalid test data"), errors.New("invalid test data")

	case "https://api.exmo.com/v1.1/currency":
		return json.Marshal(Currencies{"ADA_BTC", "ADA_USD"})

	case "https://api.exmo.com/v1.1/candles_history?symbol=ADA_BTC&resolution=30&from=1701367794&to=1701367795":
		return json.Marshal(CandlesHistory{[]Candle{{C: 1}, {C: 2}, {C: 3}}})
	case "https://api.exmo.com/v1.1/candles_history?symbol=ADA_BTC&resolution=30&from=1701289470&to=1701293070":
		return json.Marshal(CandlesHistory{[]Candle{{C: 1}, {C: 2}, {C: 3}}})
	case "https://api.exmo.com/v1.1/candles_history?symbol=ADA_BTC&resolution=30&from=1701289470&to=1701296670":
		return json.Marshal(CandlesHistory{[]Candle{{C: 2}, {C: 3}, {C: 4}}})
	case "https://api.exmo.com/v1.1/candles_history?symbol=ADA_BTC&resolution=30&from=1701289470&to=1701300270":
		return json.Marshal(CandlesHistory{[]Candle{{C: 5}, {C: 6}, {C: 7}}})

	default:
	}
	return []byte("invalid test data"), errors.New("invalid test data")
}
