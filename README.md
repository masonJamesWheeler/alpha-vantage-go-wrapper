# AlphaVantage API Wrapper for Go

This Go library provides a comprehensive client to interact with Alpha Vantage's API, allowing users to fetch time series, crypto, and indicator data seamlessly.

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Version: v1.0.0](https://img.shields.io/badge/version-v1.0.0-blue)

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Example Usage](#example-usage)
- [Documentation](https://github.com/masonJamesWheeler/alpha-vantage-go-wrapper/wiki).
- [Contribution](#contribution)
- [License](#license)
- [Contact](#contact)

## Features

The Alpha Vantage Go Wrapper offers comprehensive capabilities for financial data retrieval tailored to diverse financial data needs. Our features are outlined below:

### **Time Series**

- **Intraday**: Access real-time and historical intraday stock data.
- **Daily**: Obtain daily open, high, low, and close (OHLC) stock data.
- **Daily Adjusted**: Daily OHLC data accounting for stock splits and dividends.
- **Weekly**: Retrieve consolidated weekly stock data.
- **Weekly Adjusted**: Weekly stock data factoring in stock splits and dividends.
- **Monthly**: Aggregated monthly stock data.
- **Monthly Adjusted**: Monthly stock data inclusive of stock splits and dividends.
- **Quote Endpoint**: Capture real-time stock data for any security.

### **Cryptocurrencies**

- **Exchange Rates Trending**: Get real-time exchange rates for digital and physical currencies.
- **Intraday Premium**: Premium intraday crypto data access.
- **Daily**: Source daily crypto OHLC data.
- **Weekly**: Aggregated weekly crypto data.
- **Monthly**: Monthly crypto data insights.

### **Technical Indicators**

Dive into technical indicator values for securities over time:

- **Trend Analysis**: 
  - SMA Trending, EMA Trending, WMA, DEMA, TEMA, TRIMA, KAMA, MAMA, VWAP Premium, T3, MACD Premium, MACDEXT, STOCH Trending, STOCHF, RSI Trending, STOCHRSI, WILLR, ADX Trending, ADXR, AROON Trending, BBANDS Trending, AD Trending, OBV Trending, HT_TRENDLINE, HT_SINE, HT_TRENDMODE, HT_DCPERIOD, HT_DCPHASE, HT_PHASOR.

- **Momentum Indicators**:
  - APO, PPO, MOM, BOP, ROC, ROCR, MFI, TRIX, ULTOSC, DX, MINUS_DI, PLUS_DI, MINUS_DM, PLUS_DM.

- **Volume Indicators**:
  - CCI Trending, CMO, AROONOSC, MIDPOINT, MIDPRICE, SAR, TRANGE, ATR, NATR, ADOSC.


## Installation

To install the AlphaVantage API Wrapper, use the standard `go get`:

```bash
go get github.com:masonJamesWheeler/alpha-vantage-go-wrapper
```

## Example Usage

In the following example, we showcase how to fetch data for Cryptocurrency (Bitcoin), TimeSeries (Intraday for MSFT), and Bollinger Bands Indicator for MSFT.

To begin with, set your API key as an environment variable. This ensures security and ease of changing the key without altering the code. If you are running your program in a terminal or command line, set the environment variable like this:

```bash
export API_KEY=your_api_key
```

Then, the Go program to access the data is as follows:

```go
package main

import (
	"fmt"
	"os"
	"github.com/masonJamesWheeler/alpha-vantage-go-wrapper/models"
	"github.com/masonJamesWheeler/alpha-vantage-go-wrapper/client"
)

func main() {
	apiKey := os.Getenv("API_KEY") // Fetch the environment variable
	cli := client.NewClient(apiKey)
	
	cryptoParams := models.CryptoParams{
		Symbol: "BTC",
		Interval: "1min",
		Market: "USD",
		DataType: "json",
	}

	tsParams := models.TimeSeriesParams{
		Symbol: "MSFT",
		Interval: "1min",
		OutputSize: "compact",
		DataType: "json",
	}

	idParams := models.IndicatorParams{
		Symbol: "MSFT",
		Interval: "1min",
		TimePeriod: 60,
		SeriesType: "close",
		OutputSize: "compact",
		DataType: "json",
	}

	cryptoResponse, err := cli.GetCryptoDaily(cryptoParams)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cryptoResponse)

	tsResponse, err := cli.GetIntraday(tsParams)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tsResponse)

	idResponse, err := cli.GetBBANDS(idParams)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(idResponse)
}
```

Given the above code, if the API calls are successful, the output format you can expect is as below:

```
Daily Prices and Volumes for Digital Currency
Digital Currency: Bitcoin (BTC)
Market: United States Dollar (USD)
Last Refreshed: 2023-09-11 00:00:00
Time Zone: UTC

Time                     Open                High                Low                 Close               Volume              MarketCap           
=================================================================================================================================================
2020-12-16 00:00:00      19426.43            21560.00            19278.60            21335.52            114306.34           114306.34           

...

Intraday (1min) open, high, low, close prices and volume
Symbol: MSFT
Last Refreshed: 2023-09-08 19:59:00
Interval: 1min
Output Size: Compact
Time Zone: US/Eastern

Time                     Open           High           Low            Close          Volume         
====================================================================================================
2023-09-08 17:59:00      334.66         334.67         334.66         334.66         44             

...

Bollinger Bands (BBANDS)
Symbol: MSFT
Last Refreshed: 2023-09-08 19:59:00
Interval: 1min
Output Size: 
Time Zone: 

Time                    Real Upper Band Real Middle Band Real Lower Band
=====================================================================
2023-08-18 05:01:00              316.13         315.45         314.77

...
```

This structure provides a readable display of the fetched data. With this, users can easily comprehend and process the obtained financial metrics.
