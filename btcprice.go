// Package btcprice is a fast, simple and clean way to get the price of Bitcoin.
// The package uses the Blockchain API as a source of the price. Btcprice can be
// used with custom source of data too.
package btcprice

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Price holds the different types of prices.
type Price struct {
	Market float64 `json:"15m"`  // 15 minutes delayed market price.
	Last   float64 `json:"last"` // The most recent market price.
	Buy    float64 `json:"buy"`  // Buy price.
	Sell   float64 `json:"sell"` // Sell price.
}

// PriceCurrency holds the price of the cryptocurrency in different fiat currencies.
type PriceCurrency struct {
	Usd Price `json:"USD"` // Prices in USD.
}

// The Blockchain API provides only the price of Bitcoin.
type Blockchain struct {
	BitcoinPrice PriceCurrency
}

// NewBlockchainService returns a new Blockchain structure with up to date prices.
func NewBlockchainService() (*Blockchain, error) {
	blockchain := Blockchain{}
	err := blockchain.Update(nil)
	return &blockchain, err
}

// Returns response from the Blockchain API.
func (Blockchain) getResponse() (*http.Response, error) {

	const url = "https://blockchain.info/ticker" // The URL of the Blockchain API.

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 seconds timeout.
	}

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("btcprice", "BTC price") // Lets remote servers understand what kind of traffic they are receiving.

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	return res, nil
}

// Returns the JSON body in bytes.
func (Blockchain) fetchBytes(resBody io.ReadCloser) ([]byte, error) {

	body, readErr := ioutil.ReadAll(resBody)
	return body, readErr
}

// Update the price of Bitcoin from the Blockchain API.
// The price can be updated from a custom source too.
// Example for a normal use of Update:
// 	blockchain.BitcoinPrice.Usd.Last // 5000.4
// 	................
// 	blockchain.Update(nil) // Updates the price from the Blockchain API.
// 	blockchain.BitcoinPrice.Usd.Last // 5020.72
// Example for a custom use of Update:
//	customSource := []byte(`{
//	   "USD" : {"15m" : 7154.09, "last" : 7154.09, "buy" : 7154.09, "sell" : 7154.09, "symbol" : "$"},
//	   "AUD" : {"15m" : 9289.28, "last" : 9289.28, "buy" : 9289.28, "sell" : 9289.28, "symbol" : "$"},
//	   "BRL" : {"15m" : 23640.05, "last" : 23640.05, "buy" : 23640.05, "sell" : 23640.05, "symbol" : "R$"},
//	   .............
//	}`)
// 	blockchain.Update(customSource)
// 	blockchain.BitcoinPrice.Usd.Last // 7154.09
func (b *Blockchain) Update(manualPrice []byte) error {
	var bodyBytes []byte
	if manualPrice != nil {
		bodyBytes = manualPrice
	} else {
		var readErr error
		res, resErr := b.getResponse() // The response from the Blockchain API.
		if resErr != nil {
			return resErr
		}

		bodyBytes, readErr = b.fetchBytes(res.Body) // The JSON response body to bytes.
		if resErr != nil {
			return readErr
		}
	}
	jsonErr := json.Unmarshal(bodyBytes, &b.BitcoinPrice) // The body bytes to BitcoinPrice structure.
	return jsonErr
}
