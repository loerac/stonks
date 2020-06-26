package main

import (
    "fmt"
    //"io/ioutil"
    //"encoding/json"
)

type Watchlist struct {
    Symbol      string
    CurrPrice   float64
    Prices      []float64
    LmtPrice    float64
    Type        string
    Average     float64
    Volume      []int
    Change      float64
    ChangePercent float64
}

var jsonWatchlist string

func calcAverage(arr []float64) float64 {
    sum := 0.0
    for i := 0; i < len(arr); i++ {
        sum += arr[i]
    }

    return (sum / float64(len(arr)))
}

func (wl *Watchlist) printInfo() {
    fmt.Println("Ticker: " + wl.Symbol)
    fmt.Println("--- Price:     ", wl.Prices[len(wl.Prices) - 1])
    fmt.Println("--- Volume:    ", wl.Volume[len(wl.Volume) - 1], (wl.Volume[len(wl.Volume) - 2] - wl.Volume[len(wl.Volume) - 1]))
    fmt.Printf ("--- Change:     %.5f %.3f%%\n", wl.Change, wl.ChangePercent)
    fmt.Printf ("--- Average:    %.5f\n", wl.Average)
}

func makeJSONFormat(isStart bool) {
    if isStart {
        jsonWatchlist = "[ "
    } else {
        jsonWatchlist = jsonWatchlist[:len(jsonWatchlist)-len(",")]
        jsonWatchlist += " ]"
    }
}

func (wl *Watchlist) makeJSONData() {
    jsonWatchlist += "{"
    jsonWatchlist += " \"Symbol\": \"" + wl.Symbol + "\","
    jsonWatchlist += " \"Price\":" + fmt.Sprintf("%f", wl.CurrPrice) + ","
    jsonWatchlist += " \"Change\":" + fmt.Sprintf("%.5f", wl.Change) + ","
    jsonWatchlist += " \"Change Percent\":" + fmt.Sprintf("%.3f", wl.ChangePercent) + ","
    jsonWatchlist += " \"Volume\":" + fmt.Sprintf("%d", wl.Volume[len(wl.Volume) - 1]) + ","
    jsonWatchlist += " \"Volume Change\":"  + fmt.Sprintf("%d", wl.Volume[len(wl.Volume) - 2] - wl.Volume[len(wl.Volume) - 1]) + ","
    jsonWatchlist += " \"Average\":" + fmt.Sprintf("%f", wl.Average)
    jsonWatchlist += " },"
}

func getJSONWatchlist() string {
    return jsonWatchlist
}
