package informer

import (
	"context"
	"log"
	"time"
	"blockchain_trade/internal/models"
	"blockchain_trade/internal/tickerstore"
)

type tickerInformer struct {
	store *tickerstore.TickerStore
	tickerChannel chan models.Ticker
	cancel context.CancelFunc
	doneChannel chan string
}

func StartInformer(tickerStore *tickerstore.TickerStore) *tickerInformer {
	ctx, cancel := context.WithCancel(context.Background())
	tickerInformer := tickerInformer{
		store: tickerStore, 
		tickerChannel: make(chan models.Ticker),
		cancel: cancel,
		doneChannel: make(chan string),
	}

	go func() {
		for {
			select {
			case <- ctx.Done():
				tickerInformer.doneChannel <- "Stop infromer"
				return
			default:
				select {
				case ticker := <- tickerInformer.tickerChannel:
					log.Printf("symbol: %s | price_24h: %f | volume_24h: %f | last_trade_price: %f\n", ticker.Symbol,)
				case <- time.After(1 * time.Second):
				}
			}
		}
	}()

	return &tickerInformer
}

func (tickerInformer *tickerInformer) Log(ticker models.Ticker) {
	tickerInformer.tickerChannel <- ticker
}

func (tickerInformer *tickerInformer) StopInformer() {
	tickerInformer.cancel()

	event := <- tickerInformer.doneChannel
	log.Println(event)
}