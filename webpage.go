package main

import (
    "fmt"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

func run() {
    mux := mux.NewRouter().StrictSlash(true)
    mux.HandleFunc("/", handleWriteWatchlist).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleWriteWatchlist(w http.ResponseWriter, r *http.Request) {
    respondWithJSON(w, r, http.StatusCreated, getJSONWatchlist())
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
    response, err := json.MarshalIndent(payload, "", "")
    if nil != err {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("HTTP 500: Internal Server Error"))
        return
    }

    w.WriteHeader(code)
    w.Write(response)
}
