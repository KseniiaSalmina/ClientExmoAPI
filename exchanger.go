package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ticker         = "/ticker"
	trades         = "/trades"
	orderBook      = "/order_book"
	currency       = "/currency"
	candlesHistory = "/candles_history"
)

type CandlesHistory struct {
	Candles []Candle `json:"candles"`
}

type Candle struct {
	T int64   `json:"t"`
	O float64 `json:"o"`
	C float64 `json:"c"`
	H float64 `json:"h"`
	L float64 `json:"l"`
	V float64 `json:"v"`
}

type Currencies []string

type OrderBook map[string]OrderBookPair

type OrderBookPair struct {
	AskQuantity string     `json:"ask_quantity"`
	AskAmount   string     `json:"ask_amount"`
	AskTop      string     `json:"ask_top"`
	BidQuantity string     `json:"bid_quantity"`
	BidAmount   string     `json:"bid_amount"`
	BidTop      string     `json:"bid_top"`
	Ask         [][]string `json:"ask"`
	Bid         [][]string `json:"bid"`
}

type Ticker map[string]TickerValue

type TickerValue struct {
	BuyPrice  string `json:"buy_price"`
	SellPrice string `json:"sell_price"`
	LastTrade string `json:"last_trade"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Avg       string `json:"avg"`
	Vol       string `json:"vol"`
	VolCurr   string `json:"vol_curr"`
	Updated   int64  `json:"updated"`
}

type Trades map[string][]Pair

type Pair struct {
	TradeID  int64     `json:"trade_id"`
	Date     int64     `json:"date"`
	Type     TypeTrade `json:"type"`
	Quantity string    `json:"quantity"`
	Price    string    `json:"price"`
	Amount   string    `json:"amount"`
}

type TypeTrade string

const (
	Buy  TypeTrade = "buy"
	Sell TypeTrade = "sell"
)

type Exmo struct {
	client    *http.Client
	url       string
	isTest    bool
	requester Requester
}

func NewExmo(opts ...func(exmo *Exmo)) *Exmo {
	e := &Exmo{client: &http.Client{}, url: "https://api.exmo.com/v1.1"}
	e.requester = NewClient(e.client)
	for _, option := range opts {
		option(e)
	}
	if e.isTest {
		e.requester = &MockClient{}
	}
	return e
}

func WithClient(client *http.Client) func(exmo *Exmo) {
	return func(e *Exmo) {
		e.client = client
	}
}

func WithURL(url string) func(exmo *Exmo) {
	return func(e *Exmo) {
		e.url = url
	}
}

func Test() func(exmo *Exmo) {
	return func(e *Exmo) {
		e.isTest = true
	}
}

func (e *Exmo) GetTicker() (Ticker, error) {
	data, err := e.requester.GetRequest("POST", e.url+ticker, nil)
	if err != nil {
		return nil, fmt.Errorf("Exmo_GetTicker -> %w", err)
	}

	tickerResp := Ticker{}
	err = json.Unmarshal(data, &tickerResp)
	if err != nil {
		return nil, fmt.Errorf("Exmo_GetTicker -> %w", err)
	}
	return tickerResp, nil
}

func (e *Exmo) GetTrades(pairs ...string) (Trades, error) {
	tradesResp := Trades{}
	for _, pair := range pairs {
		data, err := e.requester.GetRequest("POST", e.url+trades, strings.NewReader(`pair=`+pair))
		if err != nil {
			return nil, fmt.Errorf("Exmo_GetTrades -> %w", err)
		}

		err = json.Unmarshal(data, &tradesResp)
		if err != nil {
			return nil, fmt.Errorf("Exmo_GetTrades -> %w", err)
		}

	}
	return tradesResp, nil
}

func (e *Exmo) GetOrderBook(limit int, pairs ...string) (OrderBook, error) {
	orderBookResp := OrderBook{}
	limitStr := strconv.Itoa(limit)
	for _, pair := range pairs {
		data, err := e.requester.GetRequest("POST", e.url+orderBook, strings.NewReader(`pair=`+pair+`&limit=`+limitStr))
		if err != nil {
			return nil, fmt.Errorf("Exmo_GetOrderBook -> %w", err)
		}

		err = json.Unmarshal(data, &orderBookResp)
		if err != nil {
			return nil, fmt.Errorf("Exmo_GetOrderBook -> %w", err)
		}
	}
	return orderBookResp, nil
}

func (e *Exmo) GetCurrencies() (Currencies, error) {
	data, err := e.requester.GetRequest("POST", e.url+currency, nil)
	if err != nil {
		return nil, fmt.Errorf("Exmo_GetCurrencies -> %w", err)
	}

	currenciesResp := Currencies{}
	err = json.Unmarshal(data, &currenciesResp)
	if err != nil {
		return nil, fmt.Errorf("Exmo_GetCurrencies -> %w", err)
	}
	return currenciesResp, nil
}

func (e *Exmo) GetCandlesHistory(pair string, limit int, start, end time.Time) (CandlesHistory, error) {
	limitStr := strconv.Itoa(limit)
	startStr, endStr := strconv.Itoa(int(start.Unix())), strconv.Itoa(int(end.Unix()))

	data, err := e.requester.GetRequest("GET", e.url+candlesHistory+"?symbol="+pair+"&resolution="+limitStr+"&from="+startStr+"&to="+endStr, nil)
	if err != nil {
		return CandlesHistory{}, fmt.Errorf("Exmo_GetCandlesHistory -> %w", err)
	}

	candlesHistoryResp := CandlesHistory{}
	err = json.Unmarshal(data, &candlesHistoryResp)
	if err != nil {
		return CandlesHistory{}, fmt.Errorf("Exmo_GetCandlesHistory -> %w", err)
	}
	return candlesHistoryResp, nil
}

func (e *Exmo) GetClosePrice(pair string, limit int, start, end time.Time) ([]float64, error) {
	candlesHistoryResp, err := e.GetCandlesHistory(pair, limit, start, end)
	if err != nil {
		return nil, fmt.Errorf("Exmo_GetClosePrice -> %w", err)
	}

	closePrices := make([]float64, 0, len(candlesHistoryResp.Candles))
	for _, closePrice := range candlesHistoryResp.Candles {
		closePrices = append(closePrices, closePrice.C)
	}

	return closePrices, nil
}
