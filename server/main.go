package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/amir-the-h/okex"
	"github.com/amir-the-h/okex/api"
	"github.com/amir-the-h/okex/requests/rest/market"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	passphrase := os.Getenv("PASS_PHRASE")

	dest := okex.NormalServer // The main API server
	ctx := context.Background()
	client, err := api.NewClient(ctx, apiKey, secretKey, passphrase, dest)
	if err != nil {
		log.Fatalln(err)
	}
	response, err := client.Rest.Market.GetCandlesticksHistory(market.GetCandlesticks{
		InstID: "BTC-USDT-SWAP",
		Bar:    okex.Bar15m,
		After:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
		Limit:  10,
	})
	if err != nil {
		log.Fatalln(err)
	}

	for _, candle := range response.Candles {
		log.Printf("Candles %+v", candle)
	}
	t := (*time.Time)(&response.Candles[len(response.Candles)-1].TS)
	response, err = client.Rest.Market.GetCandlesticksHistory(market.GetCandlesticks{
		InstID: "BTC-USDT-SWAP",
		Bar:    okex.Bar15m,
		After:  t.UnixMilli(),
		Limit:  10,
	})

	for _, candle := range response.Candles {
		log.Printf("Candles2 %+v", candle)
	}
}
