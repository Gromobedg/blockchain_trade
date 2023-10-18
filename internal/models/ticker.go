package models

type Ticker struct {
	Symbol string
	Price24h float64
	Volume24h float64
	LastTradePrice float64
}
