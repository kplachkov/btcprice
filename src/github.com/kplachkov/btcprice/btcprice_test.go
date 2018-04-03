package btcprice

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
)

// Bitcoin test price.
var BtcPriceUsdStrTest = "1.1111111"

// JSON response from the Blockchain API in bytes.
var TestJsonBytes = []byte(`{
	  "USD" : {"15m" : ` + BtcPriceUsdStrTest + `, "last" : ` + BtcPriceUsdStrTest + `, "buy" : ` + BtcPriceUsdStrTest + `, "sell" : ` + BtcPriceUsdStrTest + `, "symbol" : "$"},
	  "AUD" : {"15m" : 9289.28, "last" : 9289.28, "buy" : 9289.28, "sell" : 9289.28, "symbol" : "$"},
	  "BRL" : {"15m" : 23640.05, "last" : 23640.05, "buy" : 23640.05, "sell" : 23640.05, "symbol" : "R$"},
	  "CAD" : {"15m" : 9205.78, "last" : 9205.78, "buy" : 9205.78, "sell" : 9205.78, "symbol" : "$"},
	  "CHF" : {"15m" : 6819.59, "last" : 6819.59, "buy" : 6819.59, "sell" : 6819.59, "symbol" : "CHF"},
	  "CLP" : {"15m" : 4320355.84, "last" : 4320355.84, "buy" : 4320355.84, "sell" : 4320355.84, "symbol" : "$"},
	  "CNY" : {"15m" : 44914.1, "last" : 44914.1, "buy" : 44914.1, "sell" : 44914.1, "symbol" : "¥"},
	  "DKK" : {"15m" : 43255.64, "last" : 43255.64, "buy" : 43255.64, "sell" : 43255.64, "symbol" : "kr"},
	  "EUR" : {"15m" : 5798.39, "last" : 5798.39, "buy" : 5798.39, "sell" : 5798.39, "symbol" : "€"},
	  "GBP" : {"15m" : 5090.95, "last" : 5090.95, "buy" : 5090.95, "sell" : 5090.95, "symbol" : "£"},
	  "HKD" : {"15m" : 56151.39, "last" : 56151.39, "buy" : 56151.39, "sell" : 56151.39, "symbol" : "$"},
	  "INR" : {"15m" : 465781.43, "last" : 465781.43, "buy" : 465781.43, "sell" : 465781.43, "symbol" : "₹"},
	  "ISK" : {"15m" : 705035.71, "last" : 705035.71, "buy" : 705035.71, "sell" : 705035.71, "symbol" : "kr"},
	  "JPY" : {"15m" : 762519.36, "last" : 762519.36, "buy" : 762519.36, "sell" : 762519.36, "symbol" : "¥"},
	  "KRW" : {"15m" : 7597788.22, "last" : 7597788.22, "buy" : 7597788.22, "sell" : 7597788.22, "symbol" : "₩"},
	  "NZD" : {"15m" : 9871.73, "last" : 9871.73, "buy" : 9871.73, "sell" : 9871.73, "symbol" : "$"},
	  "PLN" : {"15m" : 24429.43, "last" : 24429.43, "buy" : 24429.43, "sell" : 24429.43, "symbol" : "zł"},
	  "RUB" : {"15m" : 410693.14, "last" : 410693.14, "buy" : 410693.14, "sell" : 410693.14, "symbol" : "RUB"},
	  "SEK" : {"15m" : 59613.93, "last" : 59613.93, "buy" : 59613.93, "sell" : 59613.93, "symbol" : "kr"},
	  "SGD" : {"15m" : 9370.27, "last" : 9370.27, "buy" : 9370.27, "sell" : 9370.27, "symbol" : "$"},
	  "THB" : {"15m" : 223157.58, "last" : 223157.58, "buy" : 223157.58, "sell" : 223157.58, "symbol" : "฿"},
	  "TWD" : {"15m" : 208329.18, "last" : 208329.18, "buy" : 208329.18, "sell" : 208329.18, "symbol" : "NT$"}
	}`)

// Test the new Blockchain structure.
func TestNewBlockchainService(t *testing.T) {
	blockchain, err := NewBlockchainService()
	if err != nil {
		t.Errorf(err.Error())
	}
	if blockchain.BitcoinPrice.Usd.Last <= 0.0 {
		t.Errorf("The latest price of Bitcoin is not set or negative.")
	}
	if blockchain.BitcoinPrice.Usd.Market <= 0.0 {
		t.Errorf("The market price of Bitcoin is not set or negative.")
	}
	if blockchain.BitcoinPrice.Usd.Sell <= 0.0 {
		t.Errorf("The sell price of Bitcoin is not set or negative.")
	}
	if blockchain.BitcoinPrice.Usd.Buy <= 0.0 {
		t.Errorf("The buy price of Bitcoin is not set or negative.")
	}
}

// Test the response of the API.
func TestGetResponse(t *testing.T) {

	blockchain, serErr := NewBlockchainService()
	if serErr != nil {
		t.Errorf(serErr.Error())
	}
	res, resErr := blockchain.getResponse()
	if resErr != nil {
		t.Errorf(resErr.Error())
	}
	statusCode := res.StatusCode
	if statusCode < 200 && statusCode > 299 {
		t.Errorf("Response status code:\ngot(%d)\nwant(%s)", statusCode, "[200 - 299]")
	}
}

// Test the fetching of the response body.
func TestFetchBytes(t *testing.T) {
	testBytes := []byte("test") // Test bytes.
	testFail := false

	blockchain, serErr := NewBlockchainService()
	if serErr != nil {
		t.Errorf(serErr.Error())
	}
	resBody := ioutil.NopCloser(bytes.NewReader(testBytes))     // The test bytes to io.ReadCloser
	resultBodyBytes, fetchErr := blockchain.fetchBytes(resBody) // The result bytes.
	if fetchErr != nil {
		t.Errorf(fetchErr.Error())
	}
	if resultBodyBytes == nil {
		testFail = true
	} else if len(resultBodyBytes) != len(testBytes) {
		testFail = true
	} else {
		// Check if the result matches the test.
		for i := range resultBodyBytes {
			if resultBodyBytes[i] != testBytes[i] {
				testFail = true
				break
			}
		}
	}
	if testFail {
		t.Errorf("Fetch response body bytes:\ngot(%s)\nwant(%s)", resultBodyBytes, testBytes)
	}
}

// Test the update of the service.
func TestUpdate(t *testing.T) {

	blockchain, serErr := NewBlockchainService()
	if serErr != nil {
		t.Errorf(serErr.Error())
	}
	// The test Bitcoin price string to float.
	btcPriceUsdFloatTest, strConvErr := strconv.ParseFloat(BtcPriceUsdStrTest, 64) // String to float64.
	if strConvErr != nil {
		log.Fatal(strConvErr)
	}

	jsonErr := blockchain.Update(TestJsonBytes)
	if jsonErr != nil {
		t.Errorf(jsonErr.Error())
	}
	// Check the latest price of Bitcoin.
	if blockchain.BitcoinPrice.Usd.Last != btcPriceUsdFloatTest {
		t.Errorf("Last price of Bitcoin:\ngot(%f)\nwant(%f)", blockchain.BitcoinPrice.Usd.Last, btcPriceUsdFloatTest)
	}
	// Check the market price of Bitcoin.
	if blockchain.BitcoinPrice.Usd.Market != btcPriceUsdFloatTest {
		t.Errorf("Market price of Bitcoin:\ngot(%f)\nwant(%f)", blockchain.BitcoinPrice.Usd.Market, btcPriceUsdFloatTest)
	}
	// Check the buy price of Bitcoin.
	if blockchain.BitcoinPrice.Usd.Buy != btcPriceUsdFloatTest {
		t.Errorf("Buy price of Bitcoin:\ngot(%f)\nwant(%f)", blockchain.BitcoinPrice.Usd.Buy, btcPriceUsdFloatTest)
	}
	// Check the sell price of Bitcoin.
	if blockchain.BitcoinPrice.Usd.Sell != btcPriceUsdFloatTest {
		t.Errorf("Sell price of Bitcoin:\ngot(%f)\nwant(%f)", blockchain.BitcoinPrice.Usd.Sell, btcPriceUsdFloatTest)
	}
}
