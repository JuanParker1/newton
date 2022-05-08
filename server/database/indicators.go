package database

type IndicatorBband struct {
	ID               uint `gorm:"primaryKey"`
	CandleId         uint `gorm:"not null"`
	Candle           Candle
	CaculatingPeriod uint    `gorm:"index, not null"`
	UpperPrice       float64 `gorm:"not null"`
	MiddlePrice      float64 `gorm:"not null"`
	LowerPrice       float64 `gorm:"not null"`
	Sigma            float64 `gorm:"not null"`
}

type IndicatorStochRsi struct {
	ID       uint `gorm:"primaryKey"`
	CandleId uint `gorm:"not null"`
	Candle   Candle
	Period   uint    `gorm:"index, not null"`
	PeriodK  uint    `gorm:"not null"`
	PeriodD  uint    `gorm:"not null"`
	FastK    float64 `gorm:"not null"`
	FastD    float64 `gorm:"not null"`
}

type IndicatorRsi struct {
	ID       uint `gorm:"primaryKey"`
	CandleId uint `gorm:"not null"`
	Candle   Candle
	Period   uint    `gorm:"index, not null"`
	Rsi      float64 `gorm:"not null"`
}
