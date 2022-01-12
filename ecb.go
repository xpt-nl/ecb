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

	"github.com/shopspring/decimal"
)

// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml
// https://www.ecb.europa.eu/stats/eurofxref/eurofxref-sdmx.xml

type Symbol string

const (
	USD Symbol = "USD"
	JPY Symbol = "JPY"
	BGN Symbol = "BGN"
	CZK Symbol = "CZK"
	DKK Symbol = "DKK"
	GBP Symbol = "GBP"
	HUF Symbol = "HUF"
	PLN Symbol = "PLN"
	RON Symbol = "RON"
	SEK Symbol = "SEK"
	CHF Symbol = "CHF"
	ISK Symbol = "ISK"
	NOK Symbol = "NOK"
	HRK Symbol = "HRK"
	RUB Symbol = "RUB"
	TRY Symbol = "TRY"
	AUD Symbol = "AUD"
	BRL Symbol = "BRL"
	CAD Symbol = "CAD"
	CNY Symbol = "CNY"
	HKD Symbol = "HKD"
	IDR Symbol = "IDR"
	ILS Symbol = "ILS"
	INR Symbol = "INR"
	KRW Symbol = "KRW"
	MXN Symbol = "MXN"
	MYR Symbol = "MYR"
	NZD Symbol = "NZD"
	PHP Symbol = "PHP"
	SGD Symbol = "SGD"
	THB Symbol = "THB"
	ZAR Symbol = "ZAR"
)

// EUR returns the rate of the Euro denominated in the currency passed in as an argument.
func EUR(symbol Symbol) (rate float64, err error) {
	response, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return 0.0, err
	}
	defer response.Body.Close()
	daily := daily{}
	if err := xml.NewDecoder(response.Body).Decode(&daily); err != nil {
		return 0.0, err
	}

	// log.Println(daily.Rates[0].Time)
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
func Rates() (rates map[Symbol]float64, err error) {
	response, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	daily := daily{}
	if err := xml.NewDecoder(response.Body).Decode(&daily); err != nil {
		return nil, err
	}

	rates = make(map[Symbol]float64)
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
	Symbol Symbol `xml:"currency,attr"`
	Rate   string `xml:"rate,attr"`
}

type Error string

func (e Error) Error() string { return string(e) }
