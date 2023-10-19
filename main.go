package main

import (
	"database/sql"
	"os"
	"log"
	"blockchain_trade/internal/sqlite"
	"blockchain_trade/internal/tickerstore"
	"blockchain_trade/internal/server"
)

func main() {
	dbAlreadyExists := isExist("./tickers.db")

	db, err := openDB("./tickers.db")
	if err != nil {
		log.Fatal(err)
	} else {
		tickersDBModel := sqlite.TickersDBModel{DB: db}
		if !dbAlreadyExists {
			tickersDBModel.Init()
			log.Println("DB started")
		}

		tickerStore := tickerstore.New()
		tickerServer := server.Start(tickerStore)
		log.Println("Rest started")

		tickerInformer := informer.StartInformer(tickerStore)
		log.Println("Informer started")

		tickerStore.Init(&tickersDBModel)
		log.Println("Init ticker store completed")

		tickersDBModel.Save(tickerStore)
		tickerInformer.Stop()
		tickerServer.Stop()
		tickersDBModel.Close()
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