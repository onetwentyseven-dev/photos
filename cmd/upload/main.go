package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := "8080"

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for path: ", r.URL.Path)
		log.Printf("%+v\n", r.Header)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{
			"Hello": "World",
		})
		if err != nil {
			log.Fatal("Failed to encode response", err)
		}
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for path: ", r.URL.Path)
		log.Printf("%+v\n", r.Header)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
		if err != nil {
			log.Fatal("Failed to encode response", err)
		}
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
