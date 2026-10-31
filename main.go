package main	

import (
    "io"
    "os"
    "fmt"
    "net/http"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
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

type RefreshResponse struct {
    success bool    `json:"success"`
    id      string  `json:"id"`
    err     string  `json:"err"`
}

type Result struct {
    price float64 `json:"price"`
    time  int64   `json:"time"`
}

type ExchangeRates struct {
    success   bool               `json:"success"`
    timestamp int64              `json:"timestamp"`
    base      string             `json:"base"`
    date      string             `json:"date"`
    rates     map[string]float64 `json:"rates"`
}

func generate_id() string {
    bytes := make([]byte, id_len)
    if _, err := rand.Read(bytes); err != nil {
        panic(err)
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
        if err != nil {
            fmt.Printf("refreshing failed: %v\n", err)
            continue
        }
        defer resp.Body.Close()

        var exchange_rates ExchangeRates
        body, _ := io.ReadAll(resp.Body)
        err2 := json.Unmarshal(body, &exchange_rates)
        if err2 != nil {
            fmt.Printf("refreshing failed: %v\n", err2)
            continue
        }
        if !exchange_rates.success {
            fmt.Println("refreshing failed: provider declined")
            continue
        }
        
        // add sql entry (id, price, time)
    }
}

func api_refresh(w http.ResponseWriter, r *http.Request) {  
    id := generate_id()
    currency := r.URL.RawQuery
    resp := RefreshResponse{success: true, id: id, err: ""}

    if (!is_valid_pair(currency)) {
        resp.success = false
        resp.err = "invalid pair"
        resp.id = ""
        
        resp_bytes, _ := json.Marshal(resp)
        fmt.Fprint(w, string(resp_bytes))
        return
    }

    trigger <- Message{id: id, currency: currency}
    resp_bytes, _ := json.Marshal(resp)
    fmt.Fprint(w, string(resp_bytes))
}

func api_getquote(w http.ResponseWriter, r *http.Request) {
    //SELECT * FROM rates WHERE id = $1 LIMIT 1;
}

func api_lastquote(w http.ResponseWriter, r *http.Request) {
    //SELECT * FROM your_table WHERE pair = $1 ORDER BY time DESC LIMIT 1;
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
