package server

import (
	"net/http"
	"context"
	"time"
	"log"
	"blockchain_trade/internal/tickerstore"
	
	"github.com/gorilla/mux"
)

type tickerServer struct {
	store *tickerstore.TickerStore
}

func Start(tickerStore *tickerstore.TickerStore) *http.Server {
	router := mux.NewRouter().StrictSlash(true)
	tickerServer := tickerServer{store: tickerStore}

	router.HandleFunc("/api/tickers", tickerServer.GetAllTickerHandler).Methods("GET")
	router.HandleFunc("/api/quit", tickerServer.StopApp).Methods("POST")

	server := &http.Server{Addr: "localhost:8080", Handler: router}
	go server.ListenAndServe()

	go func() {
        log.Printf("Server listening on %s\n", server.Addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()
	
	return server
}

func Stop(server *http.Server) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	if err := server.Shutdown(context); err != nil {
        log.Printf("Server shutdown failed: %v\n", err)
		return
    }
	log.Println("Server shutdown gracefully")
}