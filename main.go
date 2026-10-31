package main	

import (
    "net/http"
)

func refresh(w http.ResponseWriter, r *http.Request) {
}

func getquote(w http.ResponseWriter, r *http.Request) {
}

func lastquote(w http.ResponseWriter, r *http.Request) {
}

func main() {
    http.HandleFunc("/refresh", refresh)
    http.HandleFunc("/getquote", getquote)
    http.HandleFunc("/lastquote", lastquote)
    http.ListenAndServe(":8080", nil)
}
