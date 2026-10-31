package main	

import (
    "os"
    "fmt"
    "net/http"
    "crypto/rand"
    "encoding/hex"
)

const id_len int = 20
const n_pairs int = 6
const pairs = [n_pairs]string{"USD/EUR", "USD/MXN", "EUR/USD", "EUR/MXN", "MXN/USD", "MXN/EUR"}

var trigger chan message
var exchange_url string
var exchange_key string

type message struct {
    id string
    currency string
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

func is_valid_pair(pair string) bool {
    for i := range n_pairs {
        if pair == pairs[i] {
            return true
        }
    }
    return false
}

func background_refresh() {
    for {
        msg := <-trigger
        // refresh
    }
}

func api_refresh(w http.ResponseWriter, r *http.Request) {  
    id := generate_id()
    currency := r.URL.RawQuery

    if (!is_valid_pair(currency)) {
        // respond with error
    }

    trigger <- message{id: id, currency: currency}
    fmt.Fprint(w, id)
}

func api_getquote(w http.ResponseWriter, r *http.Request) {
}

func api_lastquote(w http.ResponseWriter, r *http.Request) {
}

func main() {
    exchange_url = os.Getenv("EXCHANGE_URL")
    exchange_key = os.Getenv("EXCHANGE_KEY")

    trigger = make(chan message)
    go background_refresh()
    
    http.HandleFunc("/refresh",   api_refresh)
    http.HandleFunc("/getquote",  api_getquote)
    http.HandleFunc("/lastquote", api_lastquote)
    
    http.ListenAndServe(":8080", nil)
}
