package main

import (
	"database/sql"
	"os"
	"os/signal"
    "syscall"
	"log"
	"blockchain_trade/internal/sqlite"
	"blockchain_trade/internal/tickerstore"
	"blockchain_trade/internal/server"
	"blockchain_trade/internal/informer"
)

func main() {
	dbAlreadyExists := isExist("./tickers.db")

	db, err := openDB("./tickers.db")
	if err != nil {
		log.Fatal(err)
	} else {
		tickersDB := sqlite.TickersDB{DB: db}
		if !dbAlreadyExists {
			tickersDB.Init()
			log.Println("DB started")
		}

		tickerStore := tickerstore.New()
		srv := server.Start(tickerStore)
		log.Println("Rest started")

		tickerInformer := informer.StartInformer(tickerStore)
		log.Println("Informer started")

		tickerStore.Init(&tickersDB)
		log.Println("Init ticker store completed")

		tickerGrabber := informer.StartGrabber(tickerStore, tickerInformer)
		log.Println("Grabbers started")

		signalCh := make(chan os.Signal, 1)
    	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

		log.Println("App ready")
		sig := <-signalCh
    	log.Printf("Received signal: %v\n", sig)

		tickerStore.FlushToDB()
		tickerGrabber.StopGrabber()
		tickerInformer.StopInformer()
		server.Stop(srv)
		tickersDB.Close()
	}
}

func openDB(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
} 

func isExist(dbFile string) bool {
	if _, err := os.Stat(dbFile); err == nil {
		return true
	} else {
		return false
	}
}