package main

import (
	"log/slog"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.RequestURI)
		w.Write([]byte("Hello"))
	})
	http.ListenAndServe("localhost:8080", nil)
}
