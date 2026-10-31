package main	

import (
    "os"
    "fmt"
    "net/http"
    "crypto/rand"
    "encoding/hex"
)

const id_len int = 20
var trigger string
var exchange_url string
var exchange_key string

func background_refresh() {
    for {
        _ = <-trigger
        // refresh
    }
}

func api_refresh(w http.ResponseWriter, r *http.Request) {
    id := generate_id()
    trigger <- id
}

func api_getquote(w http.ResponseWriter, r *http.Request) {
}

func api_lastquote(w http.ResponseWriter, r *http.Request) {
}

func generate_id() string {
  bytes := make([]byte, id_len)
  if _, err := rand.Read(bytes); err != nil {
      fmt.Fprintf(os.Stderr, "error: %v\n", err)
      os.Exit(1)
      return ""
  }
  return hex.EncodeToString(bytes)
}

func main() {
    val := generate_id()
    fmt.Println(val)
    
    exchange_url = os.Getenv("EXCHANGE_URL")
    exchange_key = os.Getenv("EXCHANGE_KEY")

    trigger = make(chan bool)
    go background_refresh()
    
    http.HandleFunc("/refresh",   api_refresh)
    http.HandleFunc("/getquote",  api_getquote)
    http.HandleFunc("/lastquote", api_lastquote)
    
    http.ListenAndServe(":8080", nil)
}
