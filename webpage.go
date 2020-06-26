package main

import (
    "fmt"
    "encoding/json"
    "log"
    "net/http"
    //"os"
    //"time"

    "github.com/gorilla/mux"
    //"github.com/joho/godotenv"
    //"github.com/davecgh/go-spew/spew"
)

func run() {
    mux := mux.NewRouter().StrictSlash(true)
    mux.HandleFunc("/", handleWriteWatchlist).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

/*
func makeMuxRouter() http.Handler {
    muxRouter := mux.NewRouter()
    muxRouter.HandleFunc("/", handleWriteWatchlist).Methods("POST")

    return muxRouter
}
*/

func handleWriteWatchlist(w http.ResponseWriter, r *http.Request) {
    fmt.Println("THis getting called")
    /*
    var wl Watchlist

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&wl)
    if nil != err {
        respondWithJSON(w, r, http.StatusBadRequest, r.Body)
        return
    }
    defer r.Body.Close()
    */

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
