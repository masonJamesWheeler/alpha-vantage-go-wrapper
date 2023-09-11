/*
// Package models provides types and functions for working with Alpha Vantage time series data.
//
// This file contains types and functions representing the interactions and responses 
// for time series data provided by the Alpha Vantage API.
// For more information about Alpha Vantage API, see https://www.alphavantage.co/documentation/.

Author: Mason Wheeler
*/

package models

import (
	"strings"
	"encoding/json"
	"fmt"
	"time"
	"sort"
	"strconv"
)

// TimeSeriesMetaData represents the metadata for the time series data.
type TimeSeriesMetaData struct {
    Information       string `json:"1. Information"`
    Symbol            string `json:"2. Symbol"`
    LastRefreshed     string `json:"3. Last Refreshed"`
    Interval          string `json:"4. Interval"`
    OutputSize        string `json:"5. Output Size,omitempty"` // Note: using omitempty here and on other optional fields
    TimeZone          string `json:"6. Time Zone"`
    TimePeriod float64 `json:"5. Time Period,omitempty"`
    SeriesType        string `json:"6. Series Type,omitempty"`
    VolumeFactor      string `json:"6. Volume Factor (vFactor),omitempty"`
}


// TimeSeriesParams represents the parameters for querying time series data
type TimeSeriesParams struct {
	Symbol        string
	Interval      string
	Month         interface{}
	OutputSize    interface{}
	DataType      interface{}
}

// OHLCV represents the Open, High, Low, Close, and Volume data for a given timestamp.
type OHLCV struct {
	Timestamp time.Time `json:"-"`
	Open      float64   `json:"1. open,string"`
	High      float64   `json:"2. high,string"`
	Low       float64   `json:"3. low,string"`
	Close     float64   `json:"4. close,string"`
	Volume    int       `json:"5. volume,string"`
}

// AdjustedOHLCV represents the Open, High, Low, Close, Adjusted Close, and Dividend data for a given timestamp.
type AdjustedOHLCV struct {
    OHLCV
    AdjustedClose float64 `json:"5. adjusted close,string"`
    Dividend      float64 `json:"7. dividend amount,string"`
}

// TimeSeriesIntraday represents the response for the Intraday data.
type TimeSeriesIntraday struct {
	MetaData   TimeSeriesMetaData `json:"Meta Data"`
	TimeSeries []OHLCV            `json:"-"`
}

// TimeSeriesDaily represents the response for the Daily data.
type TimeSeriesDaily struct {
    MetaData TimeSeriesMetaData           `json:"Meta Data"`
    TimeSeries []OHLCV                    `json:"-"`
}

// TimeSeriesDailyAdjusted represents the response for the Daily Adjusted data.
type TimeSeriesDailyAdjusted struct {
	MetaData TimeSeriesMetaData               `json:"Meta Data"`
	TimeSeries []AdjustedOHLCV                `json:"-"`
}

// TimeSeriesWeekly represents the response for the Weekly data.
type TimeSeriesWeekly struct {
	MetaData TimeSeriesMetaData               `json:"Meta Data"`
	TimeSeries []OHLCV                        `json:"-"`
}

// TimeSeriesWeeklyAdjusted represents the response for the Weekly Adjusted data.
type TimeSeriesWeeklyAdjusted struct {
	MetaData TimeSeriesMetaData               `json:"Meta Data"`
	TimeSeries []AdjustedOHLCV                `json:"-"`
}

// TimeSeriesMonthly represents the response for the Monthly data.
type TimeSeriesMonthly struct {
	MetaData TimeSeriesMetaData               `json:"Meta Data"`
	TimeSeries []OHLCV                        `json:"-"`
}

// TimeSeriesMonthlyAdjusted represents the response for the Monthly Adjusted data.
type TimeSeriesMonthlyAdjusted struct {
	MetaData TimeSeriesMetaData               `json:"Meta Data"`
	TimeSeries []AdjustedOHLCV                `json:"-"`
}

// Quote represents the response for the Quote Endpoint Trending.
type Quote struct {
    Symbol    string  `json:"01. symbol"`
    Open      float64 `json:"02. open,string"`
    High      float64 `json:"03. high,string"`
    Low       float64 `json:"04. low,string"`
    Price     float64 `json:"05. price,string"`
    Volume    int64   `json:"06. volume,string"`
    LatestTradingDay time.Time `json:"07. latest trading day"`
    PreviousClose    float64 `json:"08. previous close,string"`
    Change           float64 `json:"09. change,string"`
    ChangePercent    string  `json:"10. change percent"`
}

// UnmarshalJSON is a custom unmarshaler for the TimeSeriesIntraday struct.
func (t *TimeSeriesIntraday) UnmarshalJSON(data []byte) error {
    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
 	   return err
    }

	if metaData, ok := raw["Meta Data"].(map[string]interface{}); ok {
		t.MetaData.Information = metaData["1. Information"].(string)
		t.MetaData.Symbol = metaData["2. Symbol"].(string)
		t.MetaData.LastRefreshed = metaData["3. Last Refreshed"].(string)
		t.MetaData.Interval = metaData["4. Interval"].(string)
		t.MetaData.OutputSize = metaData["5. Output Size"].(string)
		t.MetaData.TimeZone = metaData["6. Time Zone"].(string)
	}

	for key, value := range raw {
		if strings.HasPrefix(key, "Time Series") {
			tsData, ok := value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected map for time series data")
			}

			for k, v := range tsData {
				timestamp, err := time.Parse("2006-01-02 15:04:05", k)
				if err != nil {
					return err
				}

				ohlcvData, err := json.Marshal(v)
				if err != nil {
					return err
				}

				var ohlcv OHLCV
				ohlcv.Timestamp = timestamp
				if err := json.Unmarshal(ohlcvData, &ohlcv); err != nil {
					return err
				}
				t.TimeSeries = append(t.TimeSeries, ohlcv)
			}
		}
	}

	// Sorting is now based only on timestamps within the TimeSeries
	sort.SliceStable(t.TimeSeries, func(i, j int) bool {
		return t.TimeSeries[i].Timestamp.Before(t.TimeSeries[j].Timestamp)
	})

	return nil
}

// UnmarshalJSON is a custom unmarshaler for the TimeSeriesDaily struct.
func (ts *TimeSeriesDaily) UnmarshalJSON(data []byte) error {
    // Define a helper struct to use the default unmarshal
    type Alias TimeSeriesDaily
    aux := &struct {
        RawTimeSeries map[string]OHLCV `json:"Time Series (Daily)"`
        *Alias
    }{
        Alias: (*Alias)(ts),
    }

    // Unmarshal the data into the helper struct
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    // Convert the irregular map into a slice
    ts.TimeSeries = make([]OHLCV, 0, len(aux.RawTimeSeries))
    for dateStr, ohlcv := range aux.RawTimeSeries {
        t, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            return err
        }
        ohlcv.Timestamp = t
        ts.TimeSeries = append(ts.TimeSeries, ohlcv)
    }

    // Sort the time series based on the timestamp
    sort.Slice(ts.TimeSeries, func(i, j int) bool {
        return ts.TimeSeries[i].Timestamp.Before(ts.TimeSeries[j].Timestamp)
    })

    return nil
}

// UnmarshalJSON is a custom unmarshaler for the TimeSeriesDailyAdjusted struct.
func (ts *TimeSeriesDailyAdjusted) UnmarshalJSON(data []byte) error {
    // Define a helper struct to use the default unmarshal
    type Alias TimeSeriesDailyAdjusted
    aux := &struct {
        RawTimeSeries map[string]AdjustedOHLCV `json:"Time Series (Daily Adjusted)"`
        *Alias
    }{
        Alias: (*Alias)(ts),
    }

    // Unmarshal the data into the helper struct
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    // Convert the irregular map into a slice
    ts.TimeSeries = make([]AdjustedOHLCV, 0, len(aux.RawTimeSeries))
    for dateStr, ohlcv := range aux.RawTimeSeries {
        t, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            return err
        }
        ohlcv.Timestamp = t
        ts.TimeSeries = append(ts.TimeSeries, ohlcv)
    }

    // Sort the time series based on the timestamp
    sort.Slice(ts.TimeSeries, func(i, j int) bool {
        return ts.TimeSeries[i].Timestamp.Before(ts.TimeSeries[j].Timestamp)
    })

    return nil
}

// UnmarshalJSON is a custom unmarshaler for the TimeSeriesWeekly struct.
func (ts *TimeSeriesWeekly) UnmarshalJSON(data []byte) error {
    // Define a helper struct to use the default unmarshal
    type Alias TimeSeriesWeekly
    aux := &struct {
        RawTimeSeries map[string]OHLCV `json:"Weekly Time Series"`
        *Alias
    }{
        Alias: (*Alias)(ts),
    }

    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    ts.TimeSeries = make([]OHLCV, 0, len(aux.RawTimeSeries))
    for dateStr, ohlcv := range aux.RawTimeSeries {
        t, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            return err
        }
        ohlcv.Timestamp = t
        ts.TimeSeries = append(ts.TimeSeries, ohlcv)
    }

    sort.Slice(ts.TimeSeries, func(i, j int) bool {
        return ts.TimeSeries[i].Timestamp.Before(ts.TimeSeries[j].Timestamp)
    })

    return nil
}

// UnmarshalJSON is a custom unmarshaler for the TimeSeriesWeeklyAdjusted struct.
func (ts *TimeSeriesWeeklyAdjusted) UnmarshalJSON(data []byte) error {
    // Define a helper struct to use the default unmarshal
    type Alias TimeSeriesWeeklyAdjusted
    aux := &struct {
        RawTimeSeries map[string]AdjustedOHLCV `json:"Weekly Adjusted Time Series"`
        *Alias
    }{
        Alias: (*Alias)(ts),
    }

    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    ts.TimeSeries = make([]AdjustedOHLCV, 0, len(aux.RawTimeSeries))
    for dateStr, ohlcv := range aux.RawTimeSeries {
        t, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            return err
        }
        ohlcv.Timestamp = t
        ts.TimeSeries = append(ts.TimeSeries, ohlcv)
    }

    sort.Slice(ts.TimeSeries, func(i, j int) bool {
        return ts.TimeSeries[i].Timestamp.Before(ts.TimeSeries[j].Timestamp)
    })

    return nil
}

// UnmarshalJSON is a custom unmarshaler for the TimeSeriesMonthly struct.
func (ts *TimeSeriesMonthly) UnmarshalJSON(data []byte) error {
	// Define a helper struct to use the default unmarshal
	type Alias TimeSeriesMonthly
	aux := &struct {
		RawTimeSeries map[string]OHLCV `json:"Monthly Time Series"`
		*Alias
	}{
		Alias: (*Alias)(ts),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ts.TimeSeries = make([]OHLCV, 0, len(aux.RawTimeSeries))
	for dateStr, ohlcv := range aux.RawTimeSeries {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}
		ohlcv.Timestamp = t
		ts.TimeSeries = append(ts.TimeSeries, ohlcv)
	}

	sort.Slice(ts.TimeSeries, func(i, j int) bool {
		return ts.TimeSeries[i].Timestamp.Before(ts.TimeSeries[j].Timestamp)
	})

	return nil
}

// UnmarshalJSON is a custom unmarshaler for the TimeSeriesMonthlyAdjusted struct.
func (ts *TimeSeriesMonthlyAdjusted) UnmarshalJSON(data []byte) error {
	// Define a helper struct to use the default unmarshal
	type Alias TimeSeriesMonthlyAdjusted
	aux := &struct {
		RawTimeSeries map[string]AdjustedOHLCV `json:"Monthly Adjusted Time Series"`
		*Alias
	}{
		Alias: (*Alias)(ts),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ts.TimeSeries = make([]AdjustedOHLCV, 0, len(aux.RawTimeSeries))
	for dateStr, ohlcv := range aux.RawTimeSeries {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}
		ohlcv.Timestamp = t
		ts.TimeSeries = append(ts.TimeSeries, ohlcv)
	}

	sort.Slice(ts.TimeSeries, func(i, j int) bool {
		return ts.TimeSeries[i].Timestamp.Before(ts.TimeSeries[j].Timestamp)
	})

	return nil
}

func (q *Quote) UnmarshalJSON(data []byte) error {
	// Define a helper struct to use the default unmarshal
	type Alias Quote
	aux := &struct {
		RawQuote map[string]string `json:"Global Quote"`
		*Alias
	}{
		Alias: (*Alias)(q),
	}

	// Unmarshal the data into the helper struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Map each value from RawQuote to its corresponding field in the Quote struct
	q.Symbol = aux.RawQuote["01. symbol"]

	open, err := strconv.ParseFloat(aux.RawQuote["02. open"], 64)
	if err != nil {
		return fmt.Errorf("error parsing 'open': %v", err)
	}
	q.Open = open

	high, err := strconv.ParseFloat(aux.RawQuote["03. high"], 64)
	if err != nil {
		return fmt.Errorf("error parsing 'high': %v", err)
	}
	q.High = high

	low, err := strconv.ParseFloat(aux.RawQuote["04. low"], 64)
	if err != nil {
		return fmt.Errorf("error parsing 'low': %v", err)
	}
	q.Low = low

	price, err := strconv.ParseFloat(aux.RawQuote["05. price"], 64)
	if err != nil {
		return fmt.Errorf("error parsing 'price': %v", err)
	}
	q.Price = price

	volume, err := strconv.ParseInt(aux.RawQuote["06. volume"], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing 'volume': %v", err)
	}
	q.Volume = volume

	latestTradingDay, err := time.Parse("2006-01-02", aux.RawQuote["07. latest trading day"])
	if err != nil {
		return fmt.Errorf("error parsing 'latest trading day': %v", err)
	}
	q.LatestTradingDay = latestTradingDay

	prevClose, err := strconv.ParseFloat(aux.RawQuote["08. previous close"], 64)
	if err != nil {
		return fmt.Errorf("error parsing 'previous close': %v", err)
	}
	q.PreviousClose = prevClose

	change, err := strconv.ParseFloat(aux.RawQuote["09. change"], 64)
	if err != nil {
		return fmt.Errorf("error parsing 'change': %v", err)
	}
	q.Change = change

	q.ChangePercent = aux.RawQuote["10. change percent"]

	return nil
}


// Length returns the count of time series data entries.
func (t *TimeSeriesIntraday) Length() int {
	return len(t.TimeSeries)
}

// Length returns the count of time series data entries.
func (t *TimeSeriesDaily) Length() int {
	return len(t.TimeSeries)
}

// Length returns the count of time series data entries.
func (t *TimeSeriesDailyAdjusted) Length() int {
	return len(t.TimeSeries)
}

// Length returns the count of time series data entries.
func (t *TimeSeriesWeekly) Length() int {
	return len(t.TimeSeries)
}

// Length returns the count of time series data entries.
func (t *TimeSeriesWeeklyAdjusted) Length() int {
	return len(t.TimeSeries)
}

// Length returns the count of time series data entries.
func (t *TimeSeriesMonthly) Length() int {
	return len(t.TimeSeries)
}

// Length returns the count of time series data entries.
func (t *Quote) Length() int {
	return 1
}

// String representation of the TimeSeriesIntraday for custom printing.
func (t TimeSeriesIntraday) String() string {
	var sb strings.Builder

	// First, print metadata
	sb.WriteString(t.MetaData.Information + "\n")
	sb.WriteString(fmt.Sprintf("Symbol: %s\n", t.MetaData.Symbol))
	sb.WriteString(fmt.Sprintf("Last Refreshed: %s\n", t.MetaData.LastRefreshed))
	sb.WriteString(fmt.Sprintf("Interval: %s\n", t.MetaData.Interval))
	sb.WriteString(fmt.Sprintf("Output Size: %s\n", t.MetaData.OutputSize))
	sb.WriteString(fmt.Sprintf("Time Zone: %s\n", t.MetaData.TimeZone))
	sb.WriteString("\n")

	// Define headers for the dataframe-style table
	headers := []string{"Time", "Open", "High", "Low", "Close", "Volume"}
	sb.WriteString(fmt.Sprintf("%-25s", headers[0]))  // Increase width for Time
	for _, header := range headers[1:] {
		sb.WriteString(fmt.Sprintf("%-15s", header)) // Left-justify each header with a width of 20
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 25 + (len(headers)-1)*15))  // Print line separator based on widths
	sb.WriteString("\n")

	// Loop through the TimeSeries slice
	for _, v := range t.TimeSeries {
		timeStr := v.Timestamp.Format("2006-01-02 15:04:05")
		sb.WriteString(fmt.Sprintf("%-25s%-15.2f%-15.2f%-15.2f%-15.2f%-15d\n", timeStr, v.Open, v.High, v.Low, v.Close, v.Volume))
	}

	return sb.String()
}

// String representation of the TimeSeriesDaily for custom printing.
func (t TimeSeriesDaily) String() string {
	var sb strings.Builder

	// First, print metadata
	sb.WriteString(t.MetaData.Information + "\n")
	sb.WriteString(fmt.Sprintf("Symbol: %s\n", t.MetaData.Symbol))
	sb.WriteString(fmt.Sprintf("Last Refreshed: %s\n", t.MetaData.LastRefreshed))
	sb.WriteString(fmt.Sprintf("Output Size: %s\n", t.MetaData.OutputSize))
	sb.WriteString(fmt.Sprintf("Time Zone: %s\n", t.MetaData.TimeZone))
	sb.WriteString("\n")

	// Define headers for the dataframe-style table
	headers := []string{"Time", "Open", "High", "Low", "Close", "Volume"}
	sb.WriteString(fmt.Sprintf("%-25s", headers[0]))  // Increase width for Time
	for _, header := range headers[1:] {
		sb.WriteString(fmt.Sprintf("%-15s", header)) // Left-justify each header with a width of 20
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 25 + (len(headers)-1)*15))  // Print line separator based on widths
	sb.WriteString("\n")

	// Loop through the TimeSeries slice
	for _, v := range t.TimeSeries {
		timeStr := v.Timestamp.Format("2006-01-02")
		sb.WriteString(fmt.Sprintf("%-25s%-15.2f%-15.2f%-15.2f%-15.2f%-15d\n", timeStr, v.Open, v.High, v.Low, v.Close, v.Volume))
	}

	return sb.String()
}

// String representation of the TimeSeriesDailyAdjusted for custom printing.
func (t TimeSeriesDailyAdjusted) String() string {
	var sb strings.Builder

	// First, print metadata
	sb.WriteString(t.MetaData.Information + "\n")
	sb.WriteString(fmt.Sprintf("Symbol: %s\n", t.MetaData.Symbol))
	sb.WriteString(fmt.Sprintf("Last Refreshed: %s\n", t.MetaData.LastRefreshed))
	sb.WriteString(fmt.Sprintf("Output Size: %s\n", t.MetaData.OutputSize))
	sb.WriteString(fmt.Sprintf("Time Zone: %s\n", t.MetaData.TimeZone))
	sb.WriteString("\n")

	// Define headers for the dataframe-style table
	headers := []string{"Time", "Open", "High", "Low", "Close", "Adjusted Close", "Volume", "Dividend"}
	sb.WriteString(fmt.Sprintf("%-25s", headers[0]))  // Increase width for Time
	for _, header := range headers[1:] {
		sb.WriteString(fmt.Sprintf("%-15s", header)) // Left-justify each header with a width of 20
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 25 + (len(headers)-1)*15))  // Print line separator based on widths
	sb.WriteString("\n")

	// Loop through the TimeSeries slice
	for _, v := range t.TimeSeries {
		timeStr := v.Timestamp.Format("2006-01-02")
		sb.WriteString(fmt.Sprintf("%-25s%-15.2f%-15.2f%-15.2f%-15.2f%-15.2f%-15d%-15.2f\n", timeStr, v.Open, v.High, v.Low, v.Close, v.AdjustedClose, v.Volume, v.Dividend))
	}

	return sb.String()
}

// String representation of the TimeSeriesWeekly for custom printing.
func (t TimeSeriesWeekly) String() string {
	var sb strings.Builder

	// First, print metadata
	sb.WriteString(t.MetaData.Information + "\n")
	sb.WriteString(fmt.Sprintf("Symbol: %s\n", t.MetaData.Symbol))
	sb.WriteString(fmt.Sprintf("Last Refreshed: %s\n", t.MetaData.LastRefreshed))
	sb.WriteString(fmt.Sprintf("Output Size: %s\n", t.MetaData.OutputSize))
	sb.WriteString(fmt.Sprintf("Time Zone: %s\n", t.MetaData.TimeZone))
	sb.WriteString("\n")

	// Define headers for the dataframe-style table
	headers := []string{"Time", "Open", "High", "Low", "Close", "Volume"}
	sb.WriteString(fmt.Sprintf("%-25s", headers[0]))  // Increase width for Time
	for _, header := range headers[1:] {
		sb.WriteString(fmt.Sprintf("%-15s", header)) // Left-justify each header with a width of 20
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 25 + (len(headers)-1)*15))  // Print line separator based on widths
	sb.WriteString("\n")

	// Loop through the TimeSeries slice
	for _, v := range t.TimeSeries {
		timeStr := v.Timestamp.Format("2006-01-02")
		sb.WriteString(fmt.Sprintf("%-25s%-15.2f%-15.2f%-15.2f%-15.2f%-15d\n", timeStr, v.Open, v.High, v.Low, v.Close, v.Volume))
	}

	return sb.String()
}

// String representation of the TimeSeriesWeeklyAdjusted for custom printing.
func (t TimeSeriesWeeklyAdjusted) String() string {
	var sb strings.Builder

	// First, print metadata
	sb.WriteString(t.MetaData.Information + "\n")
	sb.WriteString(fmt.Sprintf("Symbol: %s\n", t.MetaData.Symbol))
	sb.WriteString(fmt.Sprintf("Last Refreshed: %s\n", t.MetaData.LastRefreshed))
	sb.WriteString(fmt.Sprintf("Output Size: %s\n", t.MetaData.OutputSize))
	sb.WriteString(fmt.Sprintf("Time Zone: %s\n", t.MetaData.TimeZone))
	sb.WriteString("\n")

	// Define headers for the dataframe-style table
	headers := []string{"Time", "Open", "High", "Low", "Close", "Adjusted Close", "Volume", "Dividend"}
	sb.WriteString(fmt.Sprintf("%-25s", headers[0]))  // Increase width for Time
	for _, header := range headers[1:] {
		sb.WriteString(fmt.Sprintf("%-15s", header)) // Left-justify each header with a width of 20
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 25 + (len(headers)-1)*15))  // Print line separator based on widths
	sb.WriteString("\n")

	// Loop through the TimeSeries slice
	for _, v := range t.TimeSeries {
		timeStr := v.Timestamp.Format("2006-01-02")
		sb.WriteString(fmt.Sprintf("%-25s%-15.2f%-15.2f%-15.2f%-15.2f%-15.2f%-15d%-15.2f\n", timeStr, v.Open, v.High, v.Low, v.Close, v.AdjustedClose, v.Volume, v.Dividend))
	}

	return sb.String()
}

// String representation of the TimeSeriesMonthly for custom printing.
func (t TimeSeriesMonthly) String() string {
	var sb strings.Builder

	// First, print metadata
	sb.WriteString(t.MetaData.Information + "\n")
	sb.WriteString(fmt.Sprintf("Symbol: %s\n", t.MetaData.Symbol))
	sb.WriteString(fmt.Sprintf("Last Refreshed: %s\n", t.MetaData.LastRefreshed))
	sb.WriteString(fmt.Sprintf("Output Size: %s\n", t.MetaData.OutputSize))
	sb.WriteString(fmt.Sprintf("Time Zone: %s\n", t.MetaData.TimeZone))
	sb.WriteString("\n")

	// Define headers for the dataframe-style table
	headers := []string{"Time", "Open", "High", "Low", "Close", "Volume"}
	sb.WriteString(fmt.Sprintf("%-25s", headers[0]))  // Increase width for Time
	for _, header := range headers[1:] {
		sb.WriteString(fmt.Sprintf("%-15s", header)) // Left-justify each header with a width of 20
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 25 + (len(headers)-1)*15))  // Print line separator based on widths
	sb.WriteString("\n")

	// Loop through the TimeSeries slice
	for _, v := range t.TimeSeries {
		timeStr := v.Timestamp.Format("2006-01-02")
		sb.WriteString(fmt.Sprintf("%-25s%-15.2f%-15.2f%-15.2f%-15.2f%-15d\n", timeStr, v.Open, v.High, v.Low, v.Close, v.Volume))
	}

	return sb.String()
}

// String representation of the TimeSeriesMonthlyAdjusted for custom printing.
func (t TimeSeriesMonthlyAdjusted) String() string {
	var sb strings.Builder

	// First, print metadata
	sb.WriteString(t.MetaData.Information + "\n")
	sb.WriteString(fmt.Sprintf("Symbol: %s\n", t.MetaData.Symbol))
	sb.WriteString(fmt.Sprintf("Last Refreshed: %s\n", t.MetaData.LastRefreshed))
	sb.WriteString(fmt.Sprintf("Output Size: %s\n", t.MetaData.OutputSize))
	sb.WriteString(fmt.Sprintf("Time Zone: %s\n", t.MetaData.TimeZone))
	sb.WriteString("\n")

	// Define headers for the dataframe-style table
	headers := []string{"Time", "Open", "High", "Low", "Close", "Adjusted Close", "Volume", "Dividend"}
	sb.WriteString(fmt.Sprintf("%-25s", headers[0]))  // Increase width for Time
	for _, header := range headers[1:] {
		sb.WriteString(fmt.Sprintf("%-15s", header)) // Left-justify each header with a width of 20
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 25 + (len(headers)-1)*15))  // Print line separator based on widths
	sb.WriteString("\n")

	// Loop through the TimeSeries slice
	for _, v := range t.TimeSeries {
		timeStr := v.Timestamp.Format("2006-01-02")
		sb.WriteString(fmt.Sprintf("%-25s%-15.2f%-15.2f%-15.2f%-15.2f%-15.2f%-15d%-15.2f\n", timeStr, v.Open, v.High, v.Low, v.Close, v.AdjustedClose, v.Volume, v.Dividend))
	}

	return sb.String()
}

// String representation of the Quote for custom printing.
func (q Quote) String() string {
	var sb strings.Builder

	// First, print metadata
	sb.WriteString(fmt.Sprintf("Symbol: %s\n", q.Symbol))
	sb.WriteString(fmt.Sprintf("Open: %.2f\n", q.Open))
	sb.WriteString(fmt.Sprintf("High: %.2f\n", q.High))
	sb.WriteString(fmt.Sprintf("Low: %.2f\n", q.Low))
	sb.WriteString(fmt.Sprintf("Price: %.2f\n", q.Price))
	sb.WriteString(fmt.Sprintf("Volume: %d\n", q.Volume))
	sb.WriteString(fmt.Sprintf("Latest Trading Day: %s\n", q.LatestTradingDay.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("Previous Close: %.2f\n", q.PreviousClose))
	sb.WriteString(fmt.Sprintf("Change: %.2f\n", q.Change))
	sb.WriteString(fmt.Sprintf("Change Percent: %s\n", q.ChangePercent))

	return sb.String()
}