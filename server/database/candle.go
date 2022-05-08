package database

import (
	"time"
)

type Candle struct {
	ID             uint      `gorm:"primaryKey"`
	Ticker         string    `gorm:"index:idx_candle,priority:1, not null"`
	Duration       string    `gorm:"index:idx_candle,priority:1, not null"`
	OpenPrice      float64   `gorm:"not null"`
	ClosePrice     float64   `gorm:"not null"`
	HighPrice      float64   `gorm:"not null"`
	LowPrice       float64   `gorm:"not null"`
	Volume         float64   `gorm:"not null"`
	VolumeCurrency float64   `gorm:"not null"`
	Timestamp      time.Time `gorm:"index, not null"`
}
