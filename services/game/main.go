package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Println("game service listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
