package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connector struct {
	db *gorm.DB
}

func NewConnector(user string, password string, host string, port string, databaseName string) *Connector {
	db := connect(user, password, host, port, databaseName)
	return &Connector{db}
}

func connect(user string, password string, host string, port string, databaseName string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to conn db")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to load db")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db
}

func (connector *Connector) InsertCandles(candleEntities *[]Candle) {
	connector.db.CreateInBatches(candleEntities, len(*candleEntities))
}

func (connector *Connector) GetCandles(ticker string, duration string, from time.Time, to time.Time) []Candle {
	var candles []Candle
	connector.db.Where("ticker = ? AND duration = ? and timestamp between ? and ?", ticker, duration, from, to).Find(&candles)
	return candles
}

func (connector *Connector) GetLastCandle(ticker string, duration string) (*Candle, error) {
	var candles []Candle
	connector.db.Where("ticker = ? AND duration = ? order by timestamp desc limit 1", ticker, duration).Find(&candles)

	if len(candles) == 0 {
		return nil, errors.New("No Candles")
	} else {
		return &(candles[0]), nil
	}
}

func (connector *Connector) InsertBollinger(indicator *[]IndicatorBband) {
	CHUNK_SIZE := 1000
	chunks := 1 + (len(*indicator)-1)/CHUNK_SIZE

	for i := 0; i < chunks; i++ {
		from := CHUNK_SIZE * i
		var to int
		if CHUNK_SIZE*(i+1) > len(*indicator) {
			to = len(*indicator)
		} else {
			to = CHUNK_SIZE * (i + 1)
		}
		chunked := (*indicator)[from:to]
		connector.db.CreateInBatches(chunked, len(chunked))
	}
}

func (connector *Connector) InsertStochRsi(indicator *[]IndicatorStochRsi) {
	CHUNK_SIZE := 1000
	chunks := 1 + (len(*indicator)-1)/CHUNK_SIZE

	for i := 0; i < chunks; i++ {
		from := CHUNK_SIZE * i
		var to int
		if CHUNK_SIZE*(i+1) > len(*indicator) {
			to = len(*indicator)
		} else {
			to = CHUNK_SIZE * (i + 1)
		}
		chunked := (*indicator)[from:to]
		connector.db.CreateInBatches(chunked, len(chunked))
	}
}

func (connector *Connector) InsertRsi(indicator *[]IndicatorRsi) {
	CHUNK_SIZE := 1000
	chunks := 1 + (len(*indicator)-1)/CHUNK_SIZE

	for i := 0; i < chunks; i++ {
		from := CHUNK_SIZE * i
		var to int
		if CHUNK_SIZE*(i+1) > len(*indicator) {
			to = len(*indicator)
		} else {
			to = CHUNK_SIZE * (i + 1)
		}
		chunked := (*indicator)[from:to]
		connector.db.CreateInBatches(chunked, len(chunked))
	}
}

func (connector *Connector) Genesis() {
	connector.db.Exec("DROP TABLE IF EXISTS indicator_bbands")
	connector.db.Exec("DROP TABLE IF EXISTS indicator_stoch_rsis")
	connector.db.Exec("DROP TABLE IF EXISTS indicator_rsis")
	connector.db.Exec("DROP TABLE IF EXISTS candles")
	connector.db.Migrator().CreateTable(&Candle{})
	connector.db.Migrator().CreateTable(&IndicatorBband{})
	connector.db.Migrator().CreateTable(&IndicatorStochRsi{})
	connector.db.Migrator().CreateTable(&IndicatorRsi{})
}

func (connector *Connector) Truncate() {
	log.Println("Truncate... TO..DO...")
}
