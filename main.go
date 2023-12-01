package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"time"
)

type Indicatorer interface {
	SMA(pair string, limit, period int, from, to time.Time) ([]float64, error)
	EMA(pair string, limit, period int, from, to time.Time) ([]float64, error)
}

type Exchanger interface {
	GetTicker() (Ticker, error)
	GetTrades(pairs ...string) (Trades, error)
	GetOrderBook(limit int, pairs ...string) (OrderBook, error)
	GetCurrencies() (Currencies, error)
	GetCandlesHistory(pair string, limit int, start, end time.Time) (CandlesHistory, error)
	GetClosePrice(pair string, limit int, start, end time.Time) ([]float64, error)
}

type Requester interface {
	GetRequest(method string, url string, body io.Reader) ([]byte, error)
}

type Indicator struct {
	exchange     Exchanger
	calculateSMA func(data []float64, period int) []float64
	calculateEMA func(data []float64, period int) []float64
}

func (i *Indicator) GetDataPerPeriods(pair string, limit, period int, from, to time.Time) ([]float64, error) {
	var sum float64
	data := make([]float64, 0, period)
	onePeriodTime := to.Sub(from).Hours() / float64(period)
	tillEnd, _ := time.ParseDuration(fmt.Sprintf("%fh", onePeriodTime))
	end := from

	for j := 0; j < period; j++ {

		end = end.Add(tillEnd)
		if to.Before(end) || j == period-1 {
			end = to
		}

		dataOfOnePeriod, err := i.exchange.GetClosePrice(pair, limit, from, end)
		if err != nil {
			return nil, fmt.Errorf("Indicator_GetDataPerPeriods -> %w", err)
		}

		for _, onePrice := range dataOfOnePeriod {
			sum += onePrice
		}

		data = append(data, sum/float64(len(dataOfOnePeriod)))
		sum = 0
	}

	return data, nil
}

func (i *Indicator) SMA(pair string, limit, period int, from, to time.Time) ([]float64, error) {
	data, err := i.GetDataPerPeriods(pair, limit, period, from, to)
	if err != nil {
		return nil, fmt.Errorf("Indicator_SMA -> %w", err)
	}

	return i.calculateSMA(data, period), nil
}

func (i *Indicator) EMA(pair string, limit, period int, from, to time.Time) ([]float64, error) {
	data, err := i.GetDataPerPeriods(pair, limit, period, from, to)
	if err != nil {
		return nil, fmt.Errorf("Indicator_EMA -> %w", err)
	}

	return i.calculateEMA(data, period), nil
}

type IndicatorOption func(*Indicator)

func NewIndicator(exchange Exchanger, opts ...IndicatorOption) *Indicator {
	i := &Indicator{
		exchange:     exchange,
		calculateEMA: calculateEMA,
		calculateSMA: calculateSMA,
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func calculateSMA(data []float64, period int) []float64 {
	var sum float64
	res := make([]float64, 0, period)

	for i, price := range data {
		sum += price
		res = append(res, sum/float64(i+1))
	}

	return res
}

func WithSMA(SMA func(data []float64, period int) []float64) IndicatorOption {
	return func(i *Indicator) {
		i.calculateSMA = SMA
	}
}

func calculateEMA(data []float64, period int) []float64 {
	res := make([]float64, 0, period)
	factor, _ := decimal.NewFromInt(2).Div(decimal.NewFromInt(1 + int64(period))).Float64()

	for i, _ := range data {
		if i == 0 {
			res = append(res, data[i]*factor)
		} else {
			res = append(res, data[i]*factor+(data[i-1]*(1-factor)))
		}
	}

	return res
}

func WithEMA(EMA func(data []float64, period int) []float64) IndicatorOption {
	return func(i *Indicator) {
		i.calculateEMA = EMA
	}
}

func main() {
	var exchange Exchanger
	exchange = NewExmo()
	indicator := NewIndicator(exchange)

	sma, err := indicator.SMA("BTC_USD", 30, 5, time.Now().AddDate(0, 0, -2), time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sma)

	ema, err := indicator.EMA("BTC_USD", 30, 5, time.Now().AddDate(0, 0, -2), time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ema)
}
