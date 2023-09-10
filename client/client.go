/*
// Package client offers a comprehensive client for accessing Alpha Vantage's API.
//
// The client package has been expanded to support time series, crypto, and indicator data 
// retrieval from the Alpha Vantage API. Additionally, it comprises structs for the parameters 
// associated with each method. 
//
// Detailed example usage, including setups and explanations, can be found in our README on GitHub:
// https://github.com/masonJamesWheeler/alpha-vantage-go-wrapper/blob/main/README.md
//
// For more about the Alpha Vantage API, please see: https://www.alphavantage.co/documentation/.
//
// Author: Mason Wheeler
*/

package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/masonJamesWheeler/alpha-vantage-go-wrapper/models"
	"encoding/json"
)

const alphaVantageURL = "https://www.alphavantage.co/query"

// Client represents the Alpha Vantage client
type Client struct {
	apiKey string
}

// NewClient creates a new Alpha Vantage client
func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// getTimeSeriesData retrieves time series data based on the provided parameters.
func (c *Client) getTimeSeriesData(function string, params models.TimeSeriesParams) ([]byte, error) {
	queryParams := url.Values{}
	queryParams.Add("function", function)
	queryParams.Add("symbol", params.Symbol)
	queryParams.Add("interval", params.Interval)

	if adjStr, ok := params.Adjusted.(string); ok {
		queryParams.Add("adjusted", adjStr)
	} else if adjPtr, ok := params.Adjusted.(*string); ok {
		queryParams.Add("adjusted", *adjPtr)
	}

	if extStr, ok := params.ExtendedHours.(string); ok {
		queryParams.Add("adjusted", extStr)
	} else if extPtr, ok := params.ExtendedHours.(*string); ok {
		queryParams.Add("adjusted", *extPtr)
	}

	if monthStr, ok := params.Month.(string); ok {
		queryParams.Add("month", monthStr)
	} else if monthPtr, ok := params.Month.(*string); ok {
		queryParams.Add("month", *monthPtr)
	}

	if outputStr, ok := params.OutputSize.(string); ok {
		queryParams.Add("outputsize", outputStr)
	} else if outputPtr, ok := params.OutputSize.(*string); ok {
		queryParams.Add("outputsize", *outputPtr)
	}

	if dataTypeStr, ok := params.DataType.(string); ok {
		queryParams.Add("datatype", dataTypeStr)
	} else if dataTypePtr, ok := params.DataType.(*string); ok {
		queryParams.Add("datatype", *dataTypePtr)
	}

	queryParams.Add("apikey", c.apiKey)

	resp, err := http.Get(alphaVantageURL + "?" + queryParams.Encode())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetIndicatorData retrieves indicator data based on the provided parameters.
func (c *Client) GetIndicatorData(params models.IndicatorParams) ([]byte, error) {
	queryParams := url.Values{}
	queryParams.Add("function", params.Function)
	queryParams.Add("symbol", params.Symbol)
	queryParams.Add("interval", params.Interval)
	queryParams.Add("time_period", fmt.Sprintf("%d", params.TimePeriod))
	queryParams.Add("series_type", params.SeriesType)

	if params.DataType != "" {
		queryParams.Add("datatype", params.DataType)
	}

	if params.Month != "" {
		queryParams.Add("month", params.Month)
	}

	if params.OutputSize != "" {
		queryParams.Add("outputsize", params.OutputSize)
	}

	queryParams.Add("apikey", c.apiKey)

	resp, err := http.Get(alphaVantageURL + "?" + queryParams.Encode())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}


func (c *Client) getIndicator(indicatorName string, params models.IndicatorParams) (*models.IndicatorResponse, error) {
	// Add the function name to the params
	params.Function = indicatorName
	// Fetch the data using HTTP, similar to before.
	data, err := c.GetIndicatorData(params)
	if err != nil {
		return nil, err
	}

	var indicatorResponse models.IndicatorResponse
	if err := models.UnmarshalIndicatorJSON(&indicatorResponse, data, indicatorName); err != nil {
		return nil, err
	}

	return &indicatorResponse, nil
}

// GetCryptoExchangeRates retrieves crypto exchange rates based on the provided parameters.
func (c *Client) GetCryptoExchangeRates(params models.CryptoExchangeRateParams) ([]byte, error) {
	queryParams := url.Values{}
	queryParams.Add("function", params.Function)
	queryParams.Add("from_currency", params.FromCurrency)
	queryParams.Add("to_currency", params.ToCurrency)
	queryParams.Add("apikey", c.apiKey)

	resp, err := http.Get(alphaVantageURL + "?" + queryParams.Encode())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetCurrencyExchangeRate retrieves currency exchange rates based on the provided parameters.
func (c *Client) GetCurrencyExchangeRate(params models.CurrencyExchangeParams) (*models.CurrencyExchangeRateResponse, error) {
	queryParams := url.Values{}
	queryParams.Add("function", "CURRENCY_EXCHANGE_RATE")
	queryParams.Add("from_currency", params.FromCurrency)
	queryParams.Add("to_currency", params.ToCurrency)
	queryParams.Add("apikey", c.apiKey)

	resp, err := http.Get(alphaVantageURL + "?" + queryParams.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	exchangeRateData := &models.CurrencyExchangeRateResponse{}
	err = json.Unmarshal(data, exchangeRateData)
	if err != nil {
		return nil, err
	}

	return exchangeRateData, nil
}

// getCryptoData retrieves crypto data based on the provided parameters.
func (c *Client) getCryptoData(functionType string, params models.CryptoOHLCParams) (*models.CryptoSeriesResponse, error) {
	queryParams := url.Values{}
	queryParams.Add("function", functionType)
	queryParams.Add("symbol", params.Symbol)
	queryParams.Add("interval", params.Interval)
	queryParams.Add("market", params.Market)
	if params.OutputSize != "" {
		queryParams.Add("outputsize", params.OutputSize)
	}
	if params.DataType != "" {
		queryParams.Add("datatype", params.DataType)
	}
	queryParams.Add("apikey", c.apiKey)

	resp, err := http.Get(alphaVantageURL + "?" + queryParams.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cryptoData := &models.CryptoSeriesResponse{}
	err = models.UnmarshalCryptoJSON(cryptoData, data)
	if err != nil {
		return nil, err
	}

	return cryptoData, nil
}

// GetCryptoIntraday retrieves intraday crypto data based on the provided parameters.
func (c *Client) GetCryptoIntraday(params models.CryptoOHLCParams) (*models.CryptoSeriesResponse, error) {
	return c.getCryptoData("CRYPTO_INTRADAY", params)
}

// GetCryptoDaily retrieves daily crypto data based on the provided parameters.
func (c *Client) GetCryptoDaily(params models.CryptoOHLCParams) (*models.CryptoSeriesResponse, error) {
	return c.getCryptoData("DIGITAL_CURRENCY_DAILY", params)
}

// GetCryptoWeekly retrieves weekly crypto data based on the provided parameters.
func (c *Client) GetCryptoWeekly(params models.CryptoOHLCParams) (*models.CryptoSeriesResponse, error) {
	return c.getCryptoData("DIGITAL_CURRENCY_WEEKLY", params)
}

// GetCryptoMonthly retrieves monthly crypto data based on the provided parameters.
func (c *Client) GetCryptoMonthly(params models.CryptoOHLCParams) (*models.CryptoSeriesResponse, error) {
	return c.getCryptoData("DIGITAL_CURRENCY_MONTHLY", params)
}

// GetIntraday retrieves intraday data based on the provided parameters.
// It returns a TimeSeriesIntraday and an error if there is any.
func (c *Client) GetIntraday(params models.TimeSeriesParams) (models.TimeSeriesIntraday, error) {
	data, err := c.getTimeSeriesData("TIME_SERIES_INTRADAY", params)
	if err != nil {
		return models.TimeSeriesIntraday{}, err
	}

	var intradayData models.TimeSeriesIntraday
	err = json.Unmarshal(data, &intradayData)
	if err != nil {
		return models.TimeSeriesIntraday{}, err
	}

	return intradayData, nil
}

// GetDaily retrieves daily data based on the provided parameters.
// It returns a TimeSeriesDaily and an error if there is any.
func (c *Client) GetDaily(params models.TimeSeriesParams) (models.TimeSeriesDaily, error) {
	data, err := c.getTimeSeriesData("TIME_SERIES_DAILY", params)
	if err != nil {
		return models.TimeSeriesDaily{}, err
	}

	var dailyData models.TimeSeriesDaily
	err = json.Unmarshal(data, &dailyData)
	if err != nil {
		return models.TimeSeriesDaily{}, err
	}

	return dailyData, nil
}

// GetDailyAdjusted retrieves daily adjusted data based on the provided parameters.
// It returns a TimeSeriesDailyAdjusted and an error if there is any.
func (c *Client) GetDailyAdjusted(params models.TimeSeriesParams) (models.TimeSeriesDailyAdjusted, error) {
	data, err := c.getTimeSeriesData("TIME_SERIES_DAILY_ADJUSTED", params)
	if err != nil {
		return models.TimeSeriesDailyAdjusted{}, err
	}

	var dailyAdjustedData models.TimeSeriesDailyAdjusted
	err = json.Unmarshal(data, &dailyAdjustedData)
	if err != nil {
		return models.TimeSeriesDailyAdjusted{}, err
	}
	return dailyAdjustedData, nil
}

// GetWeekly retrieves weekly data based on the provided parameters.
// It returns a TimeSeriesWeekly and an error if there is any.
func (c *Client) GetWeekly(params models.TimeSeriesParams) (models.TimeSeriesWeekly, error) {
	data, err := c.getTimeSeriesData("TIME_SERIES_WEEKLY", params)
	if err != nil {
		return models.TimeSeriesWeekly{}, err
	}

	var weeklyData models.TimeSeriesWeekly
	err = json.Unmarshal(data, &weeklyData)
	if err != nil {
		return models.TimeSeriesWeekly{}, err
	}
	return weeklyData, nil
}

// GetWeeklyAdjusted retrieves weekly adjusted data based on the provided parameters.
// It returns a TimeSeriesWeekly and an error if there is any.
func (c *Client) GetWeeklyAdjusted(params models.TimeSeriesParams) (models.TimeSeriesWeekly, error) {
	data, err := c.getTimeSeriesData("TIME_SERIES_WEEKLY_ADJUSTED", params)
	if err != nil {
		return models.TimeSeriesWeekly{}, err
	}

	var weeklyAdjustedData models.TimeSeriesWeekly
	err = json.Unmarshal(data, &weeklyAdjustedData)
	if err != nil {
		return models.TimeSeriesWeekly{}, err
	}
	return weeklyAdjustedData, nil
}

// GetMonthly retrieves monthly data based on the provided parameters.
// It returns a TimeSeriesMonthly and an error if there is any.
func (c *Client) GetMonthly(params models.TimeSeriesParams) (models.TimeSeriesMonthly, error) {
	data, err := c.getTimeSeriesData("TIME_SERIES_MONTHLY", params)
	if err != nil {
		return models.TimeSeriesMonthly{}, err
	}

	var monthlyData models.TimeSeriesMonthly
	err = json.Unmarshal(data, &monthlyData)
	if err != nil {
		return models.TimeSeriesMonthly{}, err
	}
	return monthlyData, nil
}

// GetMonthlyAdjusted retrieves monthly adjusted data based on the provided parameters.
// It returns a TimeSeriesMonthlyAdjusted and an error if there is any.
func (c *Client) GetMonthlyAdjusted(params models.TimeSeriesParams) (models.TimeSeriesMonthlyAdjusted, error) {
	data, err := c.getTimeSeriesData("TIME_SERIES_MONTHLY_ADJUSTED", params)
	if err != nil {
		return models.TimeSeriesMonthlyAdjusted{}, err
	}

	var monthlyAdjustedData models.TimeSeriesMonthlyAdjusted
	err = json.Unmarshal(data, &monthlyAdjustedData)
	if err != nil {
		return models.TimeSeriesMonthlyAdjusted{}, err
	}
	return monthlyAdjustedData, nil
}
// GetQuoteEndpoint retrieves the quote endpoint based on the provided parameters.
// It returns a Quote and an error if there is any.
func (c *Client) GetQuoteEndpoint(params models.TimeSeriesParams) (models.Quote, error) {
	data, err := c.getTimeSeriesData("GLOBAL_QUOTE", params)
	if err != nil {
		return models.Quote{}, err
	}

	var quote models.Quote
	err = json.Unmarshal(data, &quote)
	if err != nil {
		return models.Quote{}, err
	}
	return quote, nil
}

// Client methods for retrieving indicator data

// GetSMA retrieves SMA data based on the provided parameters.
func (c *Client) GetSMA(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("SMA", params)
}

// GetEMA retrieves EMA data based on the provided parameters.
func (c *Client) GetEMA(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("EMA", params)
}
// GetWMA retrieves WMA data based on the provided parameters.
func (c *Client) GetWMA(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("WMA", params)
}

// GetDEMA retrieves DEMA data based on the provided parameters.
func (c *Client) GetDEMA(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("DEMA", params)
}

// GetTEMA retrieves TEMA data based on the provided parameters.
func (c *Client) GetTEMA(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("TEMA", params)
}

// GetTRIMA retrieves TRIMA data based on the provided parameters.
func (c *Client) GetTRIMA(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("TRIMA", params)
}

// GetKAMA retrieves KAMA data based on the provided parameters.
func (c *Client) GetKAMA(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("KAMA", params)
}

// GetMAMA retrieves MAMA data based on the provided parameters.
func (c *Client) GetMAMA(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MAMA", params)
}

// GetVWAP retrieves VWAP data based on the provided parameters.
func (c *Client) GetVWAP(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("VWAP", params)
}

// GetT3 retrieves T3 data based on the provided parameters.
func (c *Client) GetT3(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("T3", params)
}

// GetMACD retrieves MACD data based on the provided parameters.
func (c *Client) GetMACD(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MACD", params)
}

// GetMACDEXT retrieves MACDEXT data based on the provided parameters.
func (c *Client) GetMACDEXT(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MACDEXT", params)
}

// GetSTOCH retrieves STOCH data based on the provided parameters.
func (c *Client) GetSTOCH(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("STOCH", params)
}

// GetSTOCHF retrieves STOCHF data based on the provided parameters.
func (c *Client) GetSTOCHF(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("STOCHF", params)
}

// GetRSI retrieves RSI data based on the provided parameters.
func (c *Client) GetRSI(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("RSI", params)
}

// GetSTOCHRSI retrieves STOCHRSI data based on the provided parameters.
func (c *Client) GetSTOCHRSI(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("STOCHRSI", params)
}

// GetWILLR retrieves WILLR data based on the provided parameters.
func (c *Client) GetWILLR(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("WILLR", params)
}

// GetADX retrieves ADX data based on the provided parameters.
func (c *Client) GetADX(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("ADX", params)
}

// GetADXR retrieves ADXR data based on the provided parameters.
func (c *Client) GetADXR(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("ADXR", params)
}

// GetAPO retrieves APO data based on the provided parameters.
func (c *Client) GetAPO(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("APO", params)
}

// GetPPO retrieves PPO data based on the provided parameters.
func (c *Client) GetPPO(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("PPO", params)
}

// GetMOM retrieves MOM data based on the provided parameters.
func (c *Client) GetMOM(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MOM", params)
}

// GetBOP retrieves BOP data based on the provided parameters.
func (c *Client) GetBOP(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("BOP", params)
}

// GetCCI retrieves CCI data based on the provided parameters.
func (c *Client) GetCCI(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("CCI", params)
}

// GetCMO retrieves CMO data based on the provided parameters.
func (c *Client) GetCMO(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("CMO", params)
}

// GetROC retrieves ROC data based on the provided parameters.
func (c *Client) GetROC(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("ROC", params)
}

// GetROCR retrieves ROCR data based on the provided parameters.
func (c *Client) GetROCR(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("ROCR", params)
}

// GetAROON retrieves AROON data based on the provided parameters.
func (c *Client) GetAROON(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("AROON", params)
}

// GetAROONOSC retrieves AROONOSC data based on the provided parameters.
func (c *Client) GetAROONOSC(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("AROONOSC", params)
}

// GetMFI retrieves MFI data based on the provided parameters.
func (c *Client) GetMFI(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MFI", params)
}

// GetTRIX retrieves TRIX data based on the provided parameters.
func (c *Client) GetTRIX(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("TRIX", params)
}

// GetULTOSC retrieves ULTOSC data based on the provided parameters.
func (c *Client) GetULTOSC(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("ULTOSC", params)
}

// GetDX retrieves DX data based on the provided parameters.
func (c *Client) GetDX(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("DX", params)
}

// GetMINUSDI retrieves MINUSDI data based on the provided parameters.
func (c *Client) GetMINUSDI(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MINUS_DI", params)
}

// GetPLUSDI retrieves PLUSDI data based on the provided parameters.
func (c *Client) GetPLUSDI(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("PLUS_DI", params)
}

// GetMINUSDM retrieves MINUSDM data based on the provided parameters.
func (c *Client) GetMINUSDM(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MINUS_DM", params)
}

// GetPLUSDM retrieves PLUSDM data based on the provided parameters.
func (c *Client) GetPLUSDM(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("PLUS_DM", params)
}

// GetBBANDS retrieves BBANDS data based on the provided parameters.
func (c *Client) GetBBANDS(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("BBANDS", params)
}

// GetMIDPOINT retrieves MIDPOINT data based on the provided parameters.
func (c *Client) GetMIDPOINT(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MIDPOINT", params)
}

// GetMIDPRICE retrieves MIDPRICE data based on the provided parameters.
func (c *Client) GetMIDPRICE(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("MIDPRICE", params)
}

// GetSAR retrieves SAR data based on the provided parameters.
func (c *Client) GetSAR(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("SAR", params)
}

// GetTRANGE retrieves TRANGE data based on the provided parameters.
func (c *Client) GetTRANGE(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("TRANGE", params)
}

// GetATR retrieves ATR data based on the provided parameters.
func (c *Client) GetATR(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("ATR", params)
}

// GetNATR retrieves NATR data based on the provided parameters.
func (c *Client) GetNATR(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("NATR", params)
}

// GetAD retrieves AD data based on the provided parameters.
func (c *Client) GetAD(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("AD", params)
}

// GetADOSC retrieves ADOSC data based on the provided parameters.
func (c *Client) GetADOSC(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("ADOSC", params)
}

// GetOBV retrieves OBV data based on the provided parameters.
func (c *Client) GetOBV(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("OBV", params)
}

// GetHTTRENDLINE retrieves HT_TRENDLINE data based on the provided parameters.
func (c *Client) GetHTTRENDLINE(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("HT_TRENDLINE", params)
}

// GetHTSINE retrieves HT_SINE data based on the provided parameters.
func (c *Client) GetHTSINE(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("HT_SINE", params)
}

// GetHTTRENDMODE retrieves HT_TRENDMODE data based on the provided parameters.
func (c *Client) GetHTTRENDMODE(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("HT_TRENDMODE", params)
}

// GetHTDCPERIOD retrieves HT_DCPERIOD data based on the provided parameters.
func (c *Client) GetHTDCPERIOD(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("HT_DCPERIOD", params)
}

// GetHTDCPHASE retrieves HT_DCPHASE data based on the provided parameters.
func (c *Client) GetHTDCPHASE(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("HT_DCPHASE", params)
}

// GetHTPHASOR retrieves HT_PHASOR data based on the provided parameters.
func (c *Client) GetHTPHASOR(params models.IndicatorParams) (*models.IndicatorResponse, error) {
	return c.getIndicator("HT_PHASOR", params)
}