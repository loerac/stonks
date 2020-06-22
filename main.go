package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
    "time"
)

type Stocks struct {
    Ticker  string  `json:"ticker"`
    Type    string  `json:"type"`
    Price   float64 `json:"price"`
    Amount  int     `json:"amount"`
}

type API struct {
    URL string  `json:"url"`
    Key string  `json:"key"`
}
type Quote struct {
    Open    float64 `json:"o"`
    Price   float64 `json:"c"`
    High    float64 `json:"h"`
    Low     float64 `json:"l"`
    PrvClose float64 `json:"pc"`
    Timestamp float64 `json:"t"`
    Error   string  `json:"err"`
}

type Watchlist struct {
    Symbol      string
    CurrPrice   float64
    LmtPrice    float64
    Type        string
    Average     []float64
    Volume      int
    Change      float64
    ChangePercent float64
}

func calcAverage(arr []float64) float64 {
    sum := 0.0
    for i := 0; i < len(arr); i++ {
        sum += arr[i]
    }

    return (sum / float64(len(arr)))
}

func (q *Quote) PrintQuote() {
    if q.Error != "" {
        fmt.Println("Error:     ", q.Error)
    } else {
        fmt.Println("Open:      ", q.Open)
        fmt.Println("Price:     ", q.Price)
        fmt.Println("High:      ", q.High)
        fmt.Println("Low:       ", q.Low)
        fmt.Println("Prev Close: ", q.PrvClose)
        fmt.Println("Time:      ", q.Timestamp)
    }
}

func (wl *Watchlist) printInfo() {
    fmt.Println("Ticker: " + wl.Symbol)
    fmt.Println("--- Price:     ", wl.CurrPrice)
    fmt.Printf ("--- Change:     %.5f %.3f%%\n", wl.Change, wl.ChangePercent)
    fmt.Printf ("--- Average:    %.5f\n", calcAverage(wl.Average))
}

func main() {
    sapi, err := os.Open("sapi.json")
    if (nil != err) {
        fmt.Print(err)
        return
    }
    defer sapi.Close()

    sticker, err := os.Open("stickers.json")
    if (nil != err) {
        fmt.Print(err)
        return
    }
    defer sticker.Close()

    sapiVal, _ := ioutil.ReadAll(sapi)
    stickerVal, _ := ioutil.ReadAll(sticker)

    var api API
    var stocks []Stocks
    watchlist := make(map[string]Watchlist)

    json.Unmarshal([]byte(sapiVal), &api)
    json.Unmarshal([]byte(stickerVal), &stocks)
    for i, _ := range stocks {
        fmt.Println("Ticker:     ", stocks[i].Ticker)
        fmt.Println("Order Type: ", stocks[i].Type)
        fmt.Println("Price:      ", stocks[i].Price)
        fmt.Println("Amount:     ", stocks[i].Amount)
        fmt.Println("===============================")
        watchlist[stocks[i].Ticker] = Watchlist{
            stocks[i].Ticker,
            0.0,
            stocks[i].Price,
            stocks[i].Type,
            []float64{},
            0.0,
            0.0,
            0.0,
        }
    }

    for ;; {
        for key, val := range watchlist {
            var quote Quote
            url := api.URL + "quote?symbol=" + key + "&token=" + api.Key

            req, err := http.NewRequest("GET", url, nil)
            if nil != err {
                fmt.Println("Couldn't set 'GET' request for " + key + "...skipping");
                fmt.Println(err)
                continue
            }

            req.Header.Add("cache-control", "no-cache")
            res, err := http.DefaultClient.Do(req)
            if nil != err {
                fmt.Println("Couldn't get data for " + key + ": URL(" + url + ")...skipping")
                fmt.Println(err)
                continue
            }
            defer res.Body.Close()

            body, _ := ioutil.ReadAll(res.Body)
            json.Unmarshal(body, &quote)
            if quote.Error != "" {
                fmt.Println("--> Couldn't get data for " + key +  "...skipping")
                quote.PrintQuote()
                continue
            }

            val.CurrPrice = quote.Price
            val.Average = append(val.Average, quote.Price)
            val.Change = val.CurrPrice - quote.PrvClose
            val.ChangePercent = (val.Change / quote.PrvClose) * 100
            watchlist[key] = val
            val.printInfo()
        }
        time.Sleep(time.Minute)
    }
}

