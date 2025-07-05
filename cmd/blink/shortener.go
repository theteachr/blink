package main

import (
	"blink/internal/blink"
	"blink/internal/demo"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type shortener struct {
	blink.Shortener
}

func newShortener() *shortener {
	return &shortener{demo.Shortener{}}
}

func (s *shortener) redirect(w http.ResponseWriter, r *http.Request) (int, error) {
	slug := chi.URLParam(r, "slug")

	log.Println("Got slug:", slug)

	destination, err := s.Resolve(slug)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, destination.String(), http.StatusTemporaryRedirect)

	return http.StatusTemporaryRedirect, nil
}

func (s *shortener) shorten(w http.ResponseWriter, r *http.Request) (int, error) {
	payload := struct {
		Url string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, err
	}

	u, err := url.Parse(payload.Url)
	if err != nil {
		return http.StatusBadRequest, err
	}

	slug, err := s.Shorten(u)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// FIXME
	su := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
		Path:   slug,
	}

	body := struct {
		Url string `json:"url"`
	}{}
	body.Url = su.String()

	if err := json.NewEncoder(w).Encode(body); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
