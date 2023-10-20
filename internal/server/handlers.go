package server

import (
	"encoding/json"
	"net/http"
	"log"
	"syscall"
)

type user struct {
	username string
	password string
}

func (tickerServer *tickerServer) GetAllTickerHandler(writer http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all tickers at %s\n", req.URL.Path)

	allTickers := tickerServer.store.GetAllTickers()
	renderJSON(writer, allTickers)
}

func (tickerServer *tickerServer) StopApp(writer http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()

	if ok && verifyUser(user{username, password}) {
		log.Printf("handling stop app at %s\n", req.URL.Path)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	} else {
		log.Printf("handling Unauthorized %s:%s\n", username, password)
		writer.Header().Set("WWW-Authenticate", `Basic realm="api"`)
    	http.Error(writer, "Unauthorized", http.StatusUnauthorized)
	}
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

func verifyUser(user user) bool {
	var users = map[string]string{
		"nikita": "password",
	}

	password, hasUser := users[user.username]
	return (hasUser && password == user.password)
}