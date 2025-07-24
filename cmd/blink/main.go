package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func run() error {
	// TODO: Read the following values (config file, command line, ...)
	port := 8080
	base := fmt.Sprint("http://127.0.0.1:", port)

	shortener, err := newShortener(base)
	if err != nil {
		return err
	}

	r := chi.NewRouter()

	r.Post("/", handle(shortener.shorten))
	r.Get("/{slug}", handle(shortener.redirect))

	return http.ListenAndServe(fmt.Sprint(":", port), r)
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
