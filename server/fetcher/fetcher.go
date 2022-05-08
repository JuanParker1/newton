package fetcher

import (
	"log"
	"sort"
	"time"

	"github.com/TurboKang/newton/database"
	"github.com/amir-the-h/okex"
	"github.com/amir-the-h/okex/api/rest"
	"github.com/amir-the-h/okex/requests/rest/market"
	responses "github.com/amir-the-h/okex/responses/market"
	"github.com/markcheno/go-talib"
)

type Fetcher struct {
	clientRest *rest.ClientRest
	db         *database.Connector
}

func NewFetcher(rest *rest.ClientRest, db *database.Connector) *Fetcher {
	return &Fetcher{rest, db}
}

func (fetcher *Fetcher) Migrate(ticker string) {
	barSizes := []okex.BarSize{okex.Bar1m, okex.Bar5m, okex.Bar15m, okex.Bar1H, okex.Bar4H, okex.Bar12H, okex.Bar1D}

	to := time.Now()

	for _, barSize := range barSizes {
		lastCandle, err := fetcher.db.GetLastCandle(ticker, string(barSize))
		var from time.Time
		if err != nil {
			from = time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local)
		} else {
			from = (*lastCandle).Timestamp.Add(barSize.Duration())
		}
		log.Printf("Migrate %s %s starts from %v until %v", ticker, barSize, from, to)
		fetcher.MigrateSpecificBar(ticker, barSize, from, to)
		fetcher.IndicateBBand(fetcher.db, ticker, barSize, from, to)
		fetcher.IndicateRsi(fetcher.db, ticker, barSize, from, to)
		fetcher.IndicateStochRsi(fetcher.db, ticker, barSize, from, to)
	}
}

func (fetcher *Fetcher) MigrateSpecificBar(ticker string, barSize okex.BarSize, from time.Time, to time.Time) {

	fetchTime := from

	for fetchTime.Before(to) {
		log.Printf("Fetch - %v", fetchTime)
		candles := fetcher.fetchCandlesticksHistory(ticker, barSize, fetchTime)
		if fetchTime.After(to) || len(candles) == 0 {
			break
		} else {
			fetchTime = candles[len(candles)-1].Timestamp.Add(barSize.Duration())
		}
	}
}

func (fetcher *Fetcher) fetchCandlesticksHistory(ticker string, barSize okex.BarSize, fetchTime time.Time) []database.Candle {
	if fetchTime.Hour() == 0 && fetchTime.Minute() == 0 && fetchTime.Second() == 0 {
		log.Printf("GetCandleSticksHistory - %s %s %v\n", ticker, barSize, fetchTime)
	}
	limit := 100

	var response responses.Candle
	var err error

	if fetchTime.Before(time.Now().Add(okex.Bar1D.Duration() * time.Duration(-365))) {
		response, err = fetcher.clientRest.Market.GetCandlesticksHistory(market.GetCandlesticks{
			InstID: ticker,
			Limit:  int64(limit),
			Bar:    barSize,
			After:  fetchTime.Add(barSize.Duration() * time.Duration(limit)).UnixMilli(),
		})
	} else {
		response, err = fetcher.clientRest.Market.GetCandlesticksHistory(market.GetCandlesticks{
			InstID: ticker,
			Limit:  int64(limit),
			Bar:    barSize,
			Before: fetchTime.UnixMilli(),
		})
	}

	if err != nil {
		log.Fatalf("GetCandleSticksHistory failed - %s %s %v", ticker, barSize, fetchTime)
	}
	candlesFromResponse := response.Candles
	sort.Slice(candlesFromResponse, func(i, j int) bool {
		iTime := (time.Time)(candlesFromResponse[i].TS)
		jTime := (time.Time)(candlesFromResponse[j].TS)
		return iTime.Before(jTime)
	})
	var candleEntities []database.Candle

	for _, candle := range candlesFromResponse {
		candleEntities = append(candleEntities, database.Candle{
			Ticker:         ticker,
			Duration:       string(barSize),
			OpenPrice:      candle.O,
			ClosePrice:     candle.C,
			HighPrice:      candle.H,
			LowPrice:       candle.L,
			Volume:         candle.Vol,
			VolumeCurrency: candle.VolCcy,
			Timestamp:      (time.Time)(candle.TS),
		})
	}
	fetcher.db.InsertCandles(&candleEntities)
	return candleEntities
}

func (fetcher *Fetcher) IndicateBBand(db *database.Connector, ticker string, barSize okex.BarSize, start time.Time, end time.Time) {
	candles := db.GetCandles(ticker, (string)(barSize), start.Add(barSize.Duration()*-25), end)

	var closes []float64
	for _, candle := range candles {
		closes = append(closes, candle.ClosePrice)
	}
	uppers, middles, lowers := talib.BBands(closes, 20, 2, 2, 0)

	var bollingerList []database.IndicatorBband
	for i, upper := range uppers {
		candle := candles[i]
		middle := middles[i]
		lower := lowers[i]
		bband := database.IndicatorBband{
			Candle:           candle,
			CaculatingPeriod: 20,
			UpperPrice:       upper,
			MiddlePrice:      middle,
			LowerPrice:       lower,
			Sigma:            (upper - lower) / 2.0,
		}
		bollingerList = append(bollingerList, bband)
	}

	db.InsertBollinger(&bollingerList)
}

func (fetcher *Fetcher) IndicateRsi(db *database.Connector, ticker string, barSize okex.BarSize, start time.Time, end time.Time) {
	appropriateCandles := 25
	rsiCandles := 6

	candles := db.GetCandles(ticker, (string)(barSize), start.Add(barSize.Duration()*time.Duration(-appropriateCandles)), end)

	var closes []float64
	for _, candle := range candles {
		closes = append(closes, candle.ClosePrice)
	}
	rsiList := talib.Rsi(closes, rsiCandles)

	var rsiEntityList []database.IndicatorRsi
	for i, rsiValue := range rsiList {
		candle := candles[i]
		rsi := database.IndicatorRsi{
			Candle: candle,
			Period: uint(rsiCandles),
			Rsi:    rsiValue,
		}
		rsiEntityList = append(rsiEntityList, rsi)
	}

	db.InsertRsi(&rsiEntityList)
}

func (fetcher *Fetcher) IndicateStochRsi(db *database.Connector, ticker string, barSize okex.BarSize, start time.Time, end time.Time) {
	candles := db.GetCandles(ticker, (string)(barSize), start.Add(barSize.Duration()*-25), end)

	var closes []float64
	for _, candle := range candles {
		closes = append(closes, candle.ClosePrice)
	}
	fastKList, fastDList := StochRsi(closes, 14, 3, 3, 0)

	var stochRsiList []database.IndicatorStochRsi
	for i, fastK := range fastKList {
		candle := candles[i]
		fastD := fastDList[i]
		stochRsi := database.IndicatorStochRsi{
			Candle:  candle,
			Period:  14,
			PeriodK: 3,
			PeriodD: 3,
			FastK:   fastK,
			FastD:   fastD,
		}
		stochRsiList = append(stochRsiList, stochRsi)
	}

	db.InsertStochRsi(&stochRsiList)
}

func StochRsi(inReal []float64, inTimePeriod int, inFastKPeriod int, inFastDPeriod int, inFastDMAType talib.MaType) ([]float64, []float64) {

	outFastK := make([]float64, len(inReal))
	outFastD := make([]float64, len(inReal))

	lookbackSTOCHF := (inFastKPeriod - 1) + (inFastDPeriod - 1)
	lookbackTotal := inTimePeriod + lookbackSTOCHF
	startIdx := lookbackTotal
	tempRSIBuffer := talib.Rsi(inReal, inTimePeriod)
	tempk, tempd := talib.Stoch(tempRSIBuffer, tempRSIBuffer, tempRSIBuffer, inTimePeriod, inFastKPeriod, inFastDMAType, inFastDPeriod, inFastDMAType)

	for i := startIdx; i < len(inReal); i++ {
		outFastK[i] = tempk[i]
		outFastD[i] = tempd[i]
	}

	return outFastK, outFastD
}
