package main

import (
	"net/http"
)

func init() {
	http.HandleFunc("/", HealthcheckHandler)
}

func main() {
	port := "8080"
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello")); err != nil {
		w.WriteHeader(500)
	}
}
