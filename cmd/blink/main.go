package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func run() error {
	shortener := newShortener()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi!\n"))
	})

	r.Post("/", handle(shortener.shorten))
	r.Get("/{slug}", handle(shortener.redirect))

	return http.ListenAndServe(":8080", r)
}

func handle(h func(w http.ResponseWriter, r *http.Request) (int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if code, err := h(w, r); err != nil {
			payload := struct {
				Error string `json:"error"`
			}{}
			payload.Error = err.Error()

			w.WriteHeader(code)

			if err := json.NewEncoder(w).Encode(payload); err != nil {
				log.Println("failed to write error response")
			}
		}
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
