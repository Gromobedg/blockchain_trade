package tickerstore

import (
	"sync"
	"log"
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
	tickers, err := tickersDB.GetAll()
	tickerStore.tickers = tickers
	if err != nil {
		log.Fatalf("Init ticker store failed: %s\n", err)
	} 
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

func (tickerStore *TickerStore) FlushToDB() {
	if err := tickerStore.tickersDB.Flush(tickerStore.tickers); err != nil {
		log.Printf("Flush failed: %v\n", err)
	}
}

func (tickerStore *TickerStore) Save(ticker models.Ticker) bool {
	tickerStore.Lock()
	defer tickerStore.Unlock()

	currenTicker := tickerStore.tickers[ticker.Symbol]

	if currenTicker == ticker {
		return false
	} else {
		tickerStore.tickers[ticker.Symbol] = ticker
	}

	return true
}