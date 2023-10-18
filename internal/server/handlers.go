package server

import (
	"encoding/json"
	"net/http"
	"log"
)

func (tickerServer *tickerServer) GetAllTickerHandler(writer http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all tasks at %s\n", req.URL.Path)

	allTickers := tickerServer.store.GetAllTickers()
	renderJSON(writer, allTickers)
}

func renderJSON(writer http.ResponseWriter, value interface{}) {
	json, err := json.Marshal(value)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json)
}