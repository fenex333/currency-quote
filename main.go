package main	

import (
    "os"
    "net/http"
)

var exchange_url string
var exchange_key string

func background_refresh() {
}

func api_refresh(w http.ResponseWriter, r *http.Request) {
}

func api_getquote(w http.ResponseWriter, r *http.Request) {
}

func api_lastquote(w http.ResponseWriter, r *http.Request) {
}

func main() {
    exchange_url = os.Getenv("EXCHANGE_URL")
    exchange_key = os.Getenv("EXCHANGE_KEY")
    
    http.HandleFunc("/refresh",   api_refresh)
    http.HandleFunc("/getquote",  api_getquote)
    http.HandleFunc("/lastquote", api_lastquote)
    
    http.ListenAndServe(":8080", nil)
}
