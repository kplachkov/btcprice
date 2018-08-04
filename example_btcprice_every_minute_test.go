package btcprice_test

import (
	"fmt"
	"log"
	"time"

	"github.com/kplachkov/btcprice"
)

func Example() {
	blockchain, serErr := btcprice.NewBlockchainService()
	if serErr != nil {
		log.Fatal(serErr)
	}
	btcPrice := blockchain.BitcoinPrice.Usd.Last // The latest price of Bitcoin.
	// Prints the time and the latest Bitcoin price every minute if the price has changed.
	for t := range time.NewTicker(time.Second * 60).C {
		newBtcErr := blockchain.Update(nil) // Updates the price from the Blockchain API.
		if newBtcErr != nil {
			log.Fatal(newBtcErr)
		}
		newBtcPrice := blockchain.BitcoinPrice.Usd.Last
		// If the price has changed, it prints the time and up to date price.
		if btcPrice != newBtcPrice {
			fmt.Println(t.Format(time.RFC3339))
			fmt.Println("Bitcoin price:", newBtcPrice)
		}
		btcPrice = newBtcPrice
	}
}
