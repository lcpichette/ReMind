package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		hello().Render(r.Context(), w)
	})

	fmt.Printf("Listening on %v\n", "9000")
	err := http.ListenAndServe(":9000", nil)

	if err != nil {
		return
	}
}
