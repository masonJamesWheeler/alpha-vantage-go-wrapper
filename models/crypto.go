/*
// Package models provides types and functions for working with Alpha Vantage crypto data.
//
// This file contains types and functions representing the interactions and responses 
// for cryptocurrency data provided by the Alpha Vantage API.
// For more information about Alpha Vantage API, see https://www.alphavantage.co/documentation/.

Author: Mason Wheeler
*/

package models

import (
	"time"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"strconv"
)

type CryptoParams struct {
	Function   string
	Symbol     string
	Interval   string
	Market     string
	DataType   string
	OutputSize string
}

type CurrencyExchangeParams struct {
	FromCurrency string
	ToCurrency   string
}

type CryptoExchangeRateParams struct {
	Function      string
	FromCurrency  string
	ToCurrency    string
	DataType      string
}

type CurrencyExchangeRateResponse struct {
	ExchangeRateInfo ExchangeRateInfo `json:"Realtime Currency Exchange Rate"`
}

type ExchangeRateInfo struct {
	FromCurrencyCode     string `json:"1. From_Currency Code"`
	FromCurrencyName     string `json:"2. From_Currency Name"`
	ToCurrencyCode       string `json:"3. To_Currency Code"`
	ToCurrencyName       string `json:"4. To_Currency Name"`
	ExchangeRate         string `json:"5. Exchange Rate"`
	LastRefreshed        string `json:"6. Last Refreshed"`
	TimeZone             string `json:"7. Time Zone"`
	BidPrice             string `json:"8. Bid Price"`
	AskPrice             string `json:"9. Ask Price"`
}

type CryptoSeriesResponse struct {
	MetaData      CryptoMetaData
	TimeSeries    []CryptoTimeSeriesData
	IntervalLabel string
}

type CryptoMetaData struct {
	Information         string
	DigitalCurrencyCode string
	DigitalCurrencyName string
	MarketCode          string
	MarketName          string
	LastRefreshed       string
	TimeZone            string
}

type CryptoTimeSeriesData struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	MarketCap float64
}

func UnmarshalCryptoJSON(c *CryptoSeriesResponse, data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Metadata extraction
	metaData, ok := raw["Meta Data"].(map[string]interface{})
	if ok {
		c.MetaData = extractCryptoMetaData(metaData)
	}

	for tsKey, tsData := range raw {
		if strings.HasPrefix(tsKey, "Time Series") {
			c.IntervalLabel = tsKey
			timeSeriesMap := tsData.(map[string]interface{})
			for date, values := range timeSeriesMap {
				timestamp, err := time.Parse("2006-01-02", date)
				if err != nil {
					return err
				}

				valuesMap, ok := values.(map[string]interface{})
				if !ok {
					return fmt.Errorf("expected map for timestamp data")
				}

				open, _ := strconv.ParseFloat(valuesMap["1a. open (USD)"].(string), 64)
				high, _ := strconv.ParseFloat(valuesMap["2a. high (USD)"].(string), 64)
				low, _ := strconv.ParseFloat(valuesMap["3a. low (USD)"].(string), 64)
				closeVal, _ := strconv.ParseFloat(valuesMap["4a. close (USD)"].(string), 64)
				volume, _ := strconv.ParseFloat(valuesMap["5. volume"].(string), 64)
				marketCap, _ := strconv.ParseFloat(valuesMap["6. market cap (USD)"].(string), 64)

				c.TimeSeries = append(c.TimeSeries, CryptoTimeSeriesData{
					Timestamp: timestamp,
					Open:      open,
					High:      high,
					Low:       low,
					Close:     closeVal,
					Volume:    volume,
					MarketCap: marketCap,
				})
			}
		}
	}

	// Sorting based on timestamps
	sort.SliceStable(c.TimeSeries, func(a, b int) bool {
		return c.TimeSeries[a].Timestamp.Before(c.TimeSeries[b].Timestamp)
	})

	return nil
}

func extractCryptoMetaData(rawData map[string]interface{}) CryptoMetaData {
	var metaData CryptoMetaData

	for key, value := range rawData {
		switch key {
		case "1. Information":
			metaData.Information = value.(string)
		case "2. Digital Currency Code":
			metaData.DigitalCurrencyCode = value.(string)
		case "3. Digital Currency Name":
			metaData.DigitalCurrencyName = value.(string)
		case "4. Market Code":
			metaData.MarketCode = value.(string)
		case "5. Market Name":
			metaData.MarketName = value.(string)
		case "6. Last Refreshed":
			metaData.LastRefreshed = value.(string)
		case "7. Time Zone":
			metaData.TimeZone = value.(string)
		}
	}
	return metaData
}

func (c CryptoSeriesResponse) String() string {
	var sb strings.Builder

	// Print metadata
	sb.WriteString(c.MetaData.Information + "\n")
	sb.WriteString(fmt.Sprintf("Digital Currency: %s (%s)\n", c.MetaData.DigitalCurrencyName, c.MetaData.DigitalCurrencyCode))
	sb.WriteString(fmt.Sprintf("Market: %s (%s)\n", c.MetaData.MarketName, c.MetaData.MarketCode))
	sb.WriteString(fmt.Sprintf("Last Refreshed: %s\n", c.MetaData.LastRefreshed))
	sb.WriteString(fmt.Sprintf("Time Zone: %s\n", c.MetaData.TimeZone))
	sb.WriteString("\n")

	headers := []string{"Time", "Open", "High", "Low", "Close", "Volume", "MarketCap"}

	// Print headers
	for _, header := range headers {
		if header == "Time" {
			sb.WriteString(fmt.Sprintf("%-25s", header)) // Increased space for the Time column
		} else {
			sb.WriteString(fmt.Sprintf("%-20s", header))
		}
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 25 + 20*(len(headers)-1))) // Adjusting the "=" line length
	sb.WriteString("\n")

	// Loop through the TimeSeries slice
	for _, v := range c.TimeSeries {
		timeStr := v.Timestamp.Format("2006-01-02 15:04:05")
		sb.WriteString(fmt.Sprintf("%-25s%-20.2f%-20.2f%-20.2f%-20.2f%-20.2f%-20.2f", timeStr, v.Open, v.High, v.Low, v.Close, v.Volume, v.MarketCap))
		sb.WriteString("\n")
	}

	return sb.String()
}


// String function to nicely format the response for the Currency Exchange Rate API
func (r CurrencyExchangeRateResponse) String() string {
	return fmt.Sprintf(
		"From: %s (%s)\nTo: %s (%s)\nExchange Rate: %s\nLast Refreshed: %s\nTime Zone: %s\nBid Price: %s\nAsk Price: %s",
		r.ExchangeRateInfo.FromCurrencyName, r.ExchangeRateInfo.FromCurrencyCode,
		r.ExchangeRateInfo.ToCurrencyName, r.ExchangeRateInfo.ToCurrencyCode,
		r.ExchangeRateInfo.ExchangeRate,
		r.ExchangeRateInfo.LastRefreshed,
		r.ExchangeRateInfo.TimeZone,
		r.ExchangeRateInfo.BidPrice,
		r.ExchangeRateInfo.AskPrice,
	)
}
