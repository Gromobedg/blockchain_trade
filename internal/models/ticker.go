package models

type Ticker struct {
	Symbol string `json:"symbol"`
	Price24h float64 `json:"price_24h"`
	Volume24h float64 `json:"volume_24h"`
	LastTradePrice float64 `json:"last_trade_price"`
}
