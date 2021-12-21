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

type Error string

func (e Error) Error() string { return string(e) }

type Daily struct {
	Rates []Rate `xml:"Cube>Cube"`
}

type Rate struct {
	Time       string     `xml:"time,attr"`
	Currencies []Currency `xml:"Cube"`
}

type Symbol string

type Currency struct {
	Symbol Symbol `xml:"currency,attr"`
	Rate   string `xml:"rate,attr"`
}

const (
	USD Symbol = "USD"
	JPY        = "JPY"
	BGN        = "BGN"
	CZK        = "CZK"
	DKK        = "DKK"
	GBP        = "GBP"
	HUF        = "HUF"
	PLN        = "PLN"
	RON        = "RON"
	SEK        = "SEK"
	CHF        = "CHF"
	ISK        = "ISK"
	NOK        = "NOK"
	HRK        = "HRK"
	RUB        = "RUB"
	TRY        = "TRY"
	AUD        = "AUD"
	BRL        = "BRL"
	CAD        = "CAD"
	CNY        = "CNY"
	HKD        = "HKD"
	IDR        = "IDR"
	ILS        = "ILS"
	INR        = "INR"
	KRW        = "KRW"
	MXN        = "MXN"
	MYR        = "MYR"
	NZD        = "NZD"
	PHP        = "PHP"
	SGD        = "SGD"
	THB        = "THB"
	ZAR        = "ZAR"
)

// EUR returns the rate of the Euro denominated in the currency passed in as an argument.
func EUR(symbol Symbol) (rate float64, err error) {
	response, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return 0.0, err
	}
	defer response.Body.Close()
	daily := Daily{}
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
