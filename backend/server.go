package main

import (
    "fmt"
    "net/http"
)

func main() {
    fmt.Println("Starting server...")
    
    // Your server logic here
    http.HandleFunc("/", handlerFunc)
    http.ListenAndServe(":5000", nil)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
    // Handle requests
    fmt.Fprintf(w, "Hello, World!")
}

