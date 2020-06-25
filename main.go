package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
    "time"
    "strconv"
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
    Prices  []float64 `json:"c"`
    Volume  []int     `json:"v"`
    High    []float64 `json:"h"`
    Timestamp []float64 `json:"t"`
    Status  string    `json:"s"`
    Low     []float64 `json:"l"`
    Open    []float64 `json:"o"`
    Error   string    `json:"error"`
}

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
            []float64{},
            stocks[i].Price,
            stocks[i].Type,
            0.0,
            []int{},
            0.0,
            0.0,
        }
    }

    for ;; {
        loc, _ := time.LoadLocation("America/New_York")
        ts := time.Now().In(loc)
        am := time.Date(
            ts.Year(), ts.Month(), ts.Day(), 15, 30, 00, 0, time.UTC,
        )
        pm := time.Date(
            ts.Year(), ts.Month(), ts.Day(), 4, 30, 00, 0, time.UTC,
        )
        om := time.Date(
            ts.Year(), ts.Month(), ts.Day(), 9, 30, 00, 0, time.UTC,
        )
        cm := time.Date(
            ts.Year(), ts.Month(), ts.Day(), 20, 0, 00, 0, time.UTC,
        )
        cm = cm.AddDate(0, 0, 1)

        if ts.After(pm) && ts.Before(om) {
            fmt.Println("Pre Market hours...")
        } else if ts.After(om) && ts.Before(am) {
            fmt.Println("Market Hours...")
        } else if ts.After(am) && ts.Before(cm) {
            fmt.Println("After Market hours....")
        } else {
            fmt.Println("Markets are closed...")
        }

        unixFrom := strconv.FormatInt(ts.AddDate(0, 0, -1).UTC().Unix(), 10)
        unixTo := strconv.FormatInt(ts.UTC().Unix(), 10)

        for key, val := range watchlist {
            var quote Quote
            url := api.URL + "stock/candle?symbol=" + key + "&resolution=1&from=" + unixFrom + "&to=" + unixTo + "&token=" + api.Key

            req, err := http.NewRequest("GET", url, nil)
            if nil != err {
                fmt.Println("--> Couldn't set 'GET' request for " + key + "...skipping");
                fmt.Println(err)
                continue
            }

            req.Header.Add("cache-control", "no-cache")
            res, err := http.DefaultClient.Do(req)
            if nil != err {
                fmt.Println("--> Couldn't get data for " + key + ": URL(" + url + ")...skipping")
                fmt.Println(err)
                continue
            }
            defer res.Body.Close()

            body, _ := ioutil.ReadAll(res.Body)
            json.Unmarshal(body, &quote)
            if quote.Error != "" {
                fmt.Println("--> Couldn't get data for " + key +  "...skipping")
                fmt.Println("Error: " + quote.Error)
                continue
            }

            val.CurrPrice = quote.Prices[len(quote.Prices) - 1]
            val.Prices = quote.Prices
            val.Volume = quote.Volume
            val.Change = val.CurrPrice - quote.Prices[0]
            val.ChangePercent = (val.Change / quote.Prices[0]) * 100
            val.Average = calcAverage(quote.Prices)
            watchlist[key] = val
            val.printInfo()
        }
        time.Sleep(time.Minute)
    }
}

