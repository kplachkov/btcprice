<p align="center"><img src="image/Btcprice.png" alt="Btcprice"></p>

[![Codecov](https://img.shields.io/codecov/c/github/kplachkov/btcprice.svg?style=flat-square)](https://codecov.io/gh/kplachkov/btcprice)
[![Build Status](https://img.shields.io/travis/kplachkov/btcprice.svg?style=flat-square)](https://travis-ci.org/kplachkov/btcprice)
[![Go Report Card](https://goreportcard.com/badge/github.com/kplachkov/btcprice?style=flat-square)](https://goreportcard.com/report/github.com/kplachkov/btcprice)
[![GitHub release](https://img.shields.io/github/release/kplachkov/btcprice.svg?style=flat-square)](https://github.com/kplachkov/btcprice/releases)

**Btcprice** is a fast, simple and clean way to get the price of Bitcoin. The package uses the Blockchain API as a source of the price. Btcprice can be used with custom source of data too.

* [Installation](#installation)
* [Getting Started](#getting-started)
* [License](#license)


## Installation

### Install via `go get`

To install btcprice, use `go get`.

```bash
$ go get github.com/kplachkov/btcprice/
```

## Getting Started

```bash
$ godoc github.com/kplachkov/btcprice/
```

### Code Example:
```go
package main

import (
	"fmt"
	"github.com/kplachkov/btcprice"
	"log"
	"time"
)

// Prints the time and the price of Bitcoin.
func printBtcPriceAndTime(btcPrice float64, t time.Time) {
	fmt.Println(t.Format(time.RFC3339))
	fmt.Println("Bitcoin price:", btcPrice)
}

// Checks for changes in the price of Bitcoin.
func checkBtcPrice(timeDuration time.Duration) {
	blockchain, serErr := btcprice.NewBlockchainService()
	if serErr != nil {
		log.Fatal(serErr)
	}
	btcPrice := blockchain.BitcoinPrice.Usd.Last // The latest price of Bitcoin.
	printBtcPriceAndTime(btcPrice, time.Now())
	// Prints the time and the latest Bitcoin price if the price has changed.
	for t := range time.NewTicker(timeDuration).C {
		newBtcErr := blockchain.Update(nil)
		if newBtcErr != nil {
			log.Fatal(newBtcErr)
		}
		newBtcPrice := blockchain.BitcoinPrice.Usd.Last
		// If the price has changed, it prints the time and up to date price.
		if btcPrice != newBtcPrice {
			printBtcPriceAndTime(newBtcPrice, t)
		}
		btcPrice = newBtcPrice
	}
}

func main() {
	checkBtcPrice(time.Second) // Check the price of Bitcoin every second.
}

```

## License

BSD 3-Clause

Copyright (c) 2018-Present, kplachkov