package informer

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"context"
	"log"
	"time"
	"blockchain_trade/internal/tickerstore"
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
			req, _ := http.NewRequest("GET", url, nil)
			for {
				res, _ := http.DefaultClient.Do(req)
				body, _ := ioutil.ReadAll(res.Body)
				log.Println(string(body))
				res.Body.Close()
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

func (tickerGrabber *tickerGrabber) StopGrabber() {
	tickerGrabber.cancel()

	for index := 0; index < tickerGrabber.countGrabbers; index++ {
		event := <- tickerGrabber.doneChannel
		log.Println(event)
	}
}