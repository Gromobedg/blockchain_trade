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

func New() *TickerStore {
	tickerStore := TickerStore{}
	return &tickerStore
}

func (tickerStore *TickerStore) Init(tickersDBModel *sqlite.TickersDBModel) {
	tickerStore.tickers, _ = tickersDBModel.GetAll()
	tickerStore.tickersDBModel = tickersDBModel
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