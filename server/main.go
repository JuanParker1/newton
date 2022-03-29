package main

import (
	"log"
	"os"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	passPhrase := os.Getenv("PASS_PHRASE")
	api := builder.DefaultAPIBuilder.APIKey(apiKey).APISecretkey(secretKey).ApiPassphrase(passPhrase).Build(goex.OKEX) //创建现货api实例
	log.Println(api.GetTicker(goex.ETH_USDT))
}
