package main

import (
    "fmt"
    "time"

    "github.com/pkg/xattr"
)

const (
    USER    string = "user"
    API     string = ".api"
    TICKER  string = ".ticker"
    UPTODATE string = "up-to-date"
    UPDATED string = "updated"
)

func MonitorJSON() {
    go func() {
        for {
            var data []byte
            data, err := xattr.Get("stickers.json", USER + TICKER + LIST)
            if err != nil {
                fmt.Println(err)
            }

            if string(data) == "updated" {
                Setxattr("stickers.json", USER + TICKER + LIST, UPTODATE)
                fmt.Println("stickers.json", "has been updated")
            }

            time.Sleep(100 * time.Millisecond)
        }
    }()
}

func Setxattr(path, name, val string) {
    if err := xattr.Set(path, name, []byte(val)); err != nil {
        fmt.Println(err)
    }
}
