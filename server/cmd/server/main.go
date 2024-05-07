package main

import (
	"log"
	"net/http"
)

func main() {
	m := http.NewServeMux()

	m.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			return
		}
	})

	log.Println("Server is running on :8080")
	http.ListenAndServe(":8080", m)
}
