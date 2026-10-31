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
var pairs = [n_pairs]string{"USD/EUR", "USD/MXN", "EUR/USD", "EUR/MXN", "MXN/USD", "MXN/EUR"}
var trigger chan Message
var exchange_url string
var exchange_key string

type Message struct {
    id string
    currency string
}

type Result struct {
    price string `json:"price"`
    time string `json:"time"`
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
    query := exchange_url + "latest?access_key=" + exchange_key
    for {
        msg := <-trigger
        resp, err := http.Get(query + "&base=" + msg.currency[:3] + "&symbols=" + msg.currency[4:])
        defer resp.Body.Close()
        if err != nil {
            fmt.Printf("refreshing failed: %v\n", err)
            continue
        }

        // add sql entry (id, price, time)
    }
}

func api_refresh(w http.ResponseWriter, r *http.Request) {  
    id := generate_id()
    currency := r.URL.RawQuery

    if (!is_valid_pair(currency)) {
        // respond with error
    }

    trigger <- Message{id: id, currency: currency}
    fmt.Fprint(w, id) // replace with json
}

func api_getquote(w http.ResponseWriter, r *http.Request) {
    // get price & time by id
}

func api_lastquote(w http.ResponseWriter, r *http.Request) {
    // get last price & time by pair
}

func main() {
    exchange_url = os.Getenv("EXCHANGE_URL")
    exchange_key = os.Getenv("EXCHANGE_KEY")

    trigger = make(chan Message)
    go background_refresh()
    
    http.HandleFunc("/refresh",   api_refresh)
    http.HandleFunc("/getquote",  api_getquote)
    http.HandleFunc("/lastquote", api_lastquote)
    
    http.ListenAndServe(":8080", nil)
}
