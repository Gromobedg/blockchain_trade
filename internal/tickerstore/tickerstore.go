package tickerstore

import (
	"sync"
	"blockchain_trade/internal/models"
	"blockchain_trade/internal/sqlite"
)

type TickerStore struct {
	sync.Mutex
	tickers map[string]models.Ticker
	tickersDBModel *sqlite.TickersDBModel
}

func New(tickersDBModel *sqlite.TickersDBModel) *TickerStore {
	tickerStore := TickerStore{}
	tickerStore.tickers, _ = tickersDBModel.GetAll()
	tickerStore.tickersDBModel = tickersDBModel
	return &tickerStore
}

func (tickerStore *TickerStore) GetAllTickers() []models.Ticker {
	tickerStore.Lock()
	defer tickerStore.Unlock()

	allTicker := make([]models.Ticker, 0, len(tickerStore.tickers))
	for _, ticker := range tickerStore.tickers {
		allTicker = append(allTicker, ticker)
	}
	return allTicker
}