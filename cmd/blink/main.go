package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func run() error {
	shortener := newShortener()

	r := chi.NewRouter()

	r.Post("/", handle(shortener.shorten))
	r.Get("/{slug}", handle(shortener.redirect))

	return http.ListenAndServe(":8080", r)
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
