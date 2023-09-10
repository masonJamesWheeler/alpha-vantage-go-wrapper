# AlphaVantage API Wrapper for Go

This Go library provides a comprehensive client to interact with Alpha Vantage's API, allowing users to fetch time series, crypto, and indicator data seamlessly.

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Version: v1.0.0](https://img.shields.io/badge/version-v1.0.0-blue)

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Example Usage](#example-usage)
- [Documentation](#documentation)
- [Contribution](#contribution)
- [License](#license)
- [Contact](#contact)

## Features

- Simple and intuitive client functions to request data.
- Structures and models that represent Alpha Vantage's response data.
- Support for time series, crypto, and indicator data retrieval.

## Installation

To install the AlphaVantage API Wrapper, use the standard `go get`:

```bash
go get github.com/YourGitHubUsername/YourRepoName
```

## Example Usage

```go
package main

import (
	"fmt"
	"os"
	"github.com/YourGitHubUsername/YourRepoName/client"
	"github.com/YourGitHubUsername/YourRepoName/models"
)

func main() {
	apiKey := os.Getenv("API_KEY") // Set your environment variable or define it here
	cli := client.NewClient(apiKey)

	cryptoParams := models.CryptoOHLCParams{
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

	idResponse, err := cli.GetSMA(idParams)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(idResponse)
}
```
