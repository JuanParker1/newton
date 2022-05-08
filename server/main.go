package main

import (
	"context"
	"log"
	"os"

	"github.com/TurboKang/newton/database"
	"github.com/TurboKang/newton/fetcher"
	"github.com/amir-the-h/okex"
	"github.com/amir-the-h/okex/api"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	passphrase := os.Getenv("PASS_PHRASE")

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")

	dest := okex.NormalServer // The main API server
	ctx := context.Background()
	client, err := api.NewClient(ctx, apiKey, secretKey, passphrase, dest)
	if err != nil {
		log.Fatalln(err)
	}

	db := database.NewConnector(dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	db.Genesis()

	tickers := []string{"BTC-USDT-SWAP", "AVAX-USDT-SWAP", "LUNA-USDT-SWAP", "AAVE-USDT-SWAP", "ETH-USDT-SWAP"}

	fetchOperator := fetcher.NewFetcher(client.Rest, db)

	for _, ticker := range tickers {
		fetchOperator.Migrate(ticker)
	}

}
