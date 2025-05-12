/*
Package ecb provides functions to get the exchange rate between the EURO and
other currencies.

Euro foreign exchange reference rates
The reference rates are usually updated around 16:00 CET on every working day,
except on TARGET closing days. They are based on a regular daily concertation
procedure between central banks across Europe, which normally takes place
at 14:15 CET.
*/
package ecb

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

	"github.com/shopspring/decimal"
)

// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-sdmx.xml

const (
	USD = "USD"
	JPY = "JPY"
	BGN = "BGN"
	CZK = "CZK"
	DKK = "DKK"
	GBP = "GBP"
	HUF = "HUF"
	PLN = "PLN"
	RON = "RON"
	SEK = "SEK"
	CHF = "CHF"
	ISK = "ISK"
	NOK = "NOK"
	TRY = "TRY"
	AUD = "AUD"
	BRL = "BRL"
	CAD = "CAD"
	CNY = "CNY"
	HKD = "HKD"
	IDR = "IDR"
	ILS = "ILS"
	INR = "INR"
	KRW = "KRW"
	MXN = "MXN"
	MYR = "MYR"
	NZD = "NZD"
	PHP = "PHP"
	SGD = "SGD"
	THB = "THB"
	ZAR = "ZAR"
)

var Currencies = []string{
	USD,
	JPY,
	BGN,
	CZK,
	DKK,
	GBP,
	HUF,
	PLN,
	RON,
	SEK,
	CHF,
	ISK,
	NOK,
	TRY,
	AUD,
	BRL,
	CAD,
	CNY,
	HKD,
	IDR,
	ILS,
	INR,
	KRW,
	MXN,
	MYR,
	NZD,
	PHP,
	SGD,
	THB,
	ZAR,
}

// EUR returns the rate of the Euro denominated in the currency passed in as an argument.
func EUR(symbol string) (rate float64, err error) {
	daily := daily{}
	response, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return 0.0, err
	}
	defer response.Body.Close()
	if err := xml.NewDecoder(response.Body).Decode(&daily); err != nil {
		// fallback to file when response of ECB is unintelligible
		file, err := os.Open("eurofxref-daily.xml")
		if err != nil {
			return 0.0, err
		}
		defer file.Close()
		if err := xml.NewDecoder(file).Decode(&daily); err != nil {
			return 0.0, err
		}
	}

	for _, currency := range daily.Rates[0].Currencies {
		if currency.Symbol == symbol {
			cr, err := decimal.NewFromString(currency.Rate)
			rate, _ := cr.Float64()
			return rate, err
		}
	}
	return 0.0, Error(fmt.Sprint("symbol", symbol, "not found"))
}

// Rates returns all the rates of the Euro denominated in other currencies.
func EuroRates() (rates map[string]float64, err error) {
	daily := daily{}
	response, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if err := xml.NewDecoder(response.Body).Decode(&daily); err != nil {
		// fallback to file when response of ECB is unintelligible
		file, err := os.Open("eurofxref-daily.xml")
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if err := xml.NewDecoder(file).Decode(&daily); err != nil {
			return nil, err
		}
	}

	rates = make(map[string]float64)
	for _, currency := range daily.Rates[0].Currencies {
		cr, err := decimal.NewFromString(currency.Rate)
		if err != nil {
			return nil, err
		}
		rate, _ := cr.Float64()
		rates[currency.Symbol] = rate
	}
	return rates, nil
}

type daily struct {
	Rates []rate `xml:"Cube>Cube"`
}

type rate struct {
	Time       string     `xml:"time,attr"`
	Currencies []currency `xml:"Cube"`
}

type currency struct {
	Symbol string `xml:"currency,attr"`
	Rate   string `xml:"rate,attr"`
}

type Error string

func (e Error) Error() string { return string(e) }
