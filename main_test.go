package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewIndicator(t *testing.T) {
	exmo := NewExmo(Test())
	result := NewIndicator(exmo)
	assert.NotEqual(t, *result, Indicator{})
}

func TestWithSMA(t *testing.T) {
	exmo := NewExmo(Test())
	expected := []float64{1, 2, 3}
	TestSMA := func(data []float64, period int) []float64 {
		return expected
	}

	result := NewIndicator(exmo, WithSMA(TestSMA))
	returnedSMA := result.calculateSMA([]float64{}, 2)

	assert.NotEqual(t, *result, Indicator{})
	assert.NotNil(t, returnedSMA)
	assert.Equal(t, expected, returnedSMA)
}

func TestWithEMA(t *testing.T) {
	exmo := NewExmo(Test())
	expected := []float64{3, 2, 1}
	TestEMA := func(data []float64, period int) []float64 {
		return expected
	}

	result := NewIndicator(exmo, WithEMA(TestEMA))
	returnedEMA := result.calculateEMA([]float64{}, 2)

	assert.NotEqual(t, *result, Indicator{})
	assert.NotNil(t, returnedEMA)
	assert.Equal(t, expected, returnedEMA)
}

func Test_calculateSMA(t *testing.T) {
	data := []float64{1, 2, 3}
	expected := []float64{1, 1.5, 2}

	result := calculateSMA(data, len(data))

	assert.Equal(t, expected, result)
}

func Test_calculateEMA(t *testing.T) {
	data := []float64{1, 2, 3}
	expected := []float64{0.5, 1.5, 2.5}

	result := calculateEMA(data, len(data))

	assert.Equal(t, expected, result)
}

func TestIndicator_GetDataPerPeriods(t *testing.T) {
	exmo := NewExmo(Test())
	indicator := NewIndicator(exmo)

	type testData struct {
		period       int
		currencyPair string
		expected     []float64
		expectedErr  bool
	}

	testCases := []testData{
		{period: 1, currencyPair: "ADA_BTC", expected: []float64{6}, expectedErr: false},
		{period: 2, currencyPair: "BTC_USD", expected: nil, expectedErr: true},
		{period: 3, currencyPair: "ADA_BTC", expected: []float64{2, 3, 6}, expectedErr: false},
	}

	for _, tc := range testCases {
		result, err := indicator.GetDataPerPeriods(tc.currencyPair, 30, tc.period, time.Unix(1701289470, 0), time.Unix(1701300270, 0))
		if tc.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tc.expected, result)
	}
}

func TestIndicator_SMA(t *testing.T) {
	exmo := NewExmo(Test())
	indicator := NewIndicator(exmo)

	type testData struct {
		period       int
		currencyPair string
		expected     []float64
		expectedErr  bool
	}

	testCases := []testData{
		{period: 1, currencyPair: "ADA_BTC", expected: []float64{6}, expectedErr: false},
		{period: 2, currencyPair: "BTC_USD", expected: nil, expectedErr: true},
		{period: 3, currencyPair: "ADA_BTC", expected: []float64{2, 2.5, 3.6666666666666665}, expectedErr: false},
	}

	for _, tc := range testCases {
		result, err := indicator.SMA(tc.currencyPair, 30, tc.period, time.Unix(1701289470, 0), time.Unix(1701300270, 0))
		if tc.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tc.expected, result)
	}
}

func TestIndicator_EMA(t *testing.T) {
	exmo := NewExmo(Test())
	indicator := NewIndicator(exmo)

	type testData struct {
		period       int
		currencyPair string
		expected     []float64
		expectedErr  bool
	}

	testCases := []testData{
		{period: 1, currencyPair: "ADA_BTC", expected: []float64{6}, expectedErr: false},
		{period: 2, currencyPair: "BTC_USD", expected: nil, expectedErr: true},
		{period: 3, currencyPair: "ADA_BTC", expected: []float64{1, 2.5, 4.5}, expectedErr: false},
	}

	for _, tc := range testCases {
		result, err := indicator.EMA(tc.currencyPair, 30, tc.period, time.Unix(1701289470, 0), time.Unix(1701300270, 0))
		if tc.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tc.expected, result)
	}
}
