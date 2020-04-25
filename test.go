package main

import (
	"fmt"
	"net/http"
)

func testing(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("You sent a request to /hello!"))
}

func main() {
	fmt.Printf("Server started")
	http.HandleFunc("/hello", testing)
	http.ListenAndServe(":5000", nil)
}