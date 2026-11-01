package main	

import (
    "io"
    "os"
    "fmt"
    "net/http"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
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
    Success bool    `json:"success"`
    Id      string  `json:"id"`
    Err     string  `json:"err"`
}

type Result struct {
    Success bool    `json:"success"`
    Err     string  `json:"err"`
    Price   float64 `json:"price"`
    Time    int64   `json:"time"`
}

type ExchangeRates struct {
    Success   bool               `json:"success"`
    Timestamp int64              `json:"timestamp"`
    Base      string             `json:"base"`
    Date      string             `json:"date"`
    Rates     map[string]float64 `json:"rates"`
}

type RateEntry struct {
    gorm.Model
    Id   string  `gorm:"type:varchar(20);not null"`
    Pair string  `gorm:"type:varchar(7);not null"`
    Time int64   `gorm:"type:bigint(100);not null"`
    Rate float64 `gorm:"type:decimal(8,8);not null"`
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
            fmt.Printf("requesting exchange API failed: %v\n", err)
            continue
        }

        if resp.StatusCode != http.StatusOK {
            fmt.Printf("exchange API returned status: %s\n", resp.Status)
            resp.Body.Close()
            continue
        }
        
        body, err2 := io.ReadAll(resp.Body)
        resp.Body.Close()
        if err2 != nil {
            fmt.Printf("failed to read exchange response: %v\n", err2)
            continue
        }

        var exchange_rates ExchangeRates
        err3 := json.Unmarshal(body, &exchange_rates)
        if err3 != nil {
            fmt.Printf("failed to unmarshal exchange responce: %v\n", err3)
            continue
        }
        
        if !exchange_rates.Success {
            fmt.Println("refreshing failed: provider declined")
            continue
        }

        entry := RateEntry{
            Id: msg.id,
            Pair: msg.currency,
            Time: exchange_rates.timestamp,
            Rate: exchange_rates.Rates[msg.currency[4:]]}
        
        db.Create(&entry)
    }
}

func api_refresh(w http.ResponseWriter, r *http.Request) {  
    var resp RefreshResponse
    id := generate_id()
    currency := r.URL.RawQuery
    
    if (is_valid_pair(currency)) {
        resp.Success = true
        resp.Err = ""
        resp.Id = id
        trigger <- Message{id: id, currency: currency}
    } else {
        resp.Success = false
        resp.Err = "invalid pair"
        resp.Id = ""
    }

    resp_bytes, _ := json.Marshal(resp)
    fmt.Fprint(w, string(resp_bytes))
}

func api_getquote(w http.ResponseWriter, r *http.Request) {
    var entry RatesEntry
    var result Result
    
    id = r.URL.RawQuery
    err := db.First(&entry, "Id = ?", id).Error
    
    if err != nil {
        result.Success = false
        result.Err = "failed fetching database entry"
        result.Price = 0.0
        result.Time = 0
    } else {
        result.Success = true
        result.Err = ""
        result.Price = entry.Rate
        result.Time = entry.Time
    }
    
    resp_bytes, _ := json.Marshal(result)
    fmt.Fprint(w, string(resp_bytes))
}

func api_lastquote(w http.ResponseWriter, r *http.Request) {
    var entry RatesEntry
    var result Result
    
    pair = r.URL.RawQuery
    err := db.Where("Pair = ?", pair).Order("Time DESC").First(&entry).Error
    
    if err != nil {
        result.Success = false
        result.Err = "failed fetching database entry"
        result.Price = 0.0
        result.Time = 0
    } else {
        result.Success = true
        result.Err = ""
        result.Price = entry.Rate
        result.Time = entry.Time
    }
    
    resp_bytes, _ := json.Marshal(result)
    fmt.Fprint(w, string(resp_bytes))
}

func main() {
    exchange_url = os.Getenv("EXCHANGE_URL")
    exchange_key = os.Getenv("EXCHANGE_KEY")

    dsn := "host=localhost user=pg password=pg dbname=pg port=5432"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }
    db.AutoMigrate(&RateEntry{})

    trigger = make(chan Message, 100)
    go background_refresh()
    
    http.HandleFunc("/refresh",   api_refresh)
    http.HandleFunc("/getquote",  api_getquote)
    http.HandleFunc("/lastquote", api_lastquote)
    
    http.ListenAndServe(":8080", nil)
}
