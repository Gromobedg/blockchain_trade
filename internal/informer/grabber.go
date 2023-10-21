package informer

import (
	"fmt"
	"net/http"
	"context"
	"log"
	"time"
	"encoding/json"
	"runtime"
	"blockchain_trade/internal/tickerstore"
	"blockchain_trade/internal/models"
	"blockchain_trade/internal/logerror"
)

type tickerGrabber struct {
	store *tickerstore.TickerStore
	informer *tickerInformer
	cancel context.CancelFunc
	doneChannel chan string
	countGrabbers int
}

func StartGrabber(tickerStore *tickerstore.TickerStore, tickerInformer *tickerInformer) *tickerGrabber {
	ctx, cancel := context.WithCancel(context.Background())
	keys := tickerStore.GetAllKeys()
	tickerGrabber := tickerGrabber{
		store: tickerStore, 
		informer: tickerInformer,
		cancel: cancel,
		doneChannel: make(chan string),
		countGrabbers: len(keys),
	}

	for _, symbol := range keys {
		go func(symbol string) {
			url := fmt.Sprintf("https://api.blockchain.com/v3/exchange/tickers/%s", symbol)
			req, err := http.NewRequest("GET", url, nil)
			tickerGrabber.catchErr(err, symbol)
			for {
				resp, err := http.DefaultClient.Do(req)
				tickerGrabber.catchErr(err, symbol)
				
				var ticker models.Ticker
				err = json.NewDecoder(resp.Body).Decode(&ticker)
				resp.Body.Close()
				tickerGrabber.catchErr(err, symbol)

				if tickerStore.Save(ticker) {
					tickerInformer.Log(ticker)
				}

				select {
				case <- ctx.Done():
					tickerGrabber.doneChannel <- ("Stop grabber " + symbol)
					return
				case <- time.After(60 * time.Second):
				}
			}
		}(symbol)
	}

	return &tickerGrabber
}

func (tickerGrabber *tickerGrabber) catchErr(err interface{}, symbol string) {
	if err != nil {
		logerror.Log(err)
		tickerGrabber.doneChannel <- ("Stop grabber " + symbol)
		runtime.Goexit()
	}
}

func (tickerGrabber *tickerGrabber) StopGrabber() {
	tickerGrabber.cancel()

	for index := 0; index < tickerGrabber.countGrabbers; index++ {
		event := <- tickerGrabber.doneChannel
		log.Println(event)
	}
}