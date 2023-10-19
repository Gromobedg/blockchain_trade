package server

import (
	"net/http"
	"blockchain_trade/internal/tickerstore"
	
	"github.com/gorilla/mux"
)

type tickerServer struct {
	store *tickerstore.TickerStore
}

func Start(tickerStore *tickerstore.TickerStore) *tickerServer {
	router := mux.NewRouter().StrictSlash(true)
	server := tickerServer{store: tickerStore}

	router.HandleFunc("/api/tickers", server.GetAllTickerHandler).Methods("GET")
	router.HandleFunc("/api/quit", server.StopApp).Methods("POST")

	go http.ListenAndServe("localhost:8080", router)
	
	return &server
}