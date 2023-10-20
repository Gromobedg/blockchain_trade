package tickerstore

import (
	"sync"
	"blockchain_trade/internal/models"
	"blockchain_trade/internal/sqlite"
)

type TickerStore struct {
	sync.Mutex
	tickers map[string]models.Ticker
	tickersDB *sqlite.TickersDB
}

func New() *TickerStore {
	tickerStore := TickerStore{}
	return &tickerStore
}

func (tickerStore *TickerStore) Init(tickersDB *sqlite.TickersDB) {
	tickerStore.tickers, _ = tickersDB.GetAll()
	tickerStore.tickersDB = tickersDB
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

func (tickerStore *TickerStore) GetAllKeys() []string {
	tickerStore.Lock()
	defer tickerStore.Unlock()

	allkeys := make([]string, 0, len(tickerStore.tickers))
	for key, _ := range tickerStore.tickers {
		allkeys = append(allkeys, key)
	}
	return allkeys
}

func (tickerStore *TickerStore) FlushToDB() error {
	return tickerStore.tickersDB.Flush(tickerStore.tickers)
}

func (tickerStore *TickerStore) Save(ticker models.Ticker) {
	tickerStore.Lock()
	defer tickerStore.Unlock()

	tickerStore.tickers[ticker.Symbol] = ticker
}