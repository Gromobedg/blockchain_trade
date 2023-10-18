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
		tickersModel := sqlite.TickersDBModel{DB: db}
		if !dbAlreadyExists {
			tickersModel.Init()
			log.Println("Init db completed")
		}

		tickerStore := tickerstore.New(&tickersModel)
		server.StartServer(tickerStore)

		tickersModel.Close()
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