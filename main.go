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

    Setxattr("stickers.json", USER + TICKER + LIST, UPTODATE)
    MonitorJSON()

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


    go func() { run() }()
    done := make(chan bool)
    for ;; {
        makeJSONFormat(true)
        NYUTC, _ := time.LoadLocation("America/New_York")
        ts := time.Now().In(NYUTC)
        pm := time.Date(
            ts.Year(), ts.Month(), ts.Day(), 4, 00, 00, 00, NYUTC,
        )
        om := time.Date(
            ts.Year(), ts.Month(), ts.Day(), 9, 30, 00, 00, NYUTC,
        )
        am := time.Date(
            ts.Year(), ts.Month(), ts.Day(), 15, 00, 00, 00, NYUTC,
        )
        cm := time.Date(
            ts.Year(), ts.Month(), ts.Day(), 20, 00, 00, 00, NYUTC,
        )

        if ts.After(pm) && ts.Before(om) {
            fmt.Println("============== Pre Market hours ==============")
        } else if ts.After(om) && ts.Before(am) {
            fmt.Println("================ Market Hours ================")
        } else if ts.After(am) && ts.Before(cm) {
            fmt.Println("============= After Market hours =============")
        } else {
            fmt.Println("============= Markets are closed =============")
            pm = pm.AddDate(0, 0, 1)
            if 6 >= int(pm.Weekday()) {
                pm = pm.AddDate(0, 0, 8 - int(pm.Weekday()))
            }

            diff := pm.Sub(ts)
            fmt.Println("Doing one check, and then sleeping until premarket hours:", pm)
            go func() {
                <-done
                time.Sleep(diff)
                done <- true
            }()
        }

        unixFrom := strconv.FormatInt(am.AddDate(0, 0, -1).UTC().Unix(), 10)
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
                fmt.Println("--> Error: " + quote.Error)
                continue
            } else if 0 == len(quote.Prices) {
                fmt.Println("--> Couldn't get data for " + key +  "...skipping")
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
            val.makeJSONData()
        }
        makeJSONFormat(false)
        go func() { done <- true }()
        time.Sleep(time.Minute)
        <-done
    }
}

