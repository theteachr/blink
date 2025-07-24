package main

import (
	"blink/internal/blink"
	"blink/internal/demo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type shortener struct {
	blink.Shortener
	base *url.URL
}

func newShortener(base string) (*shortener, error) {
	base_, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	return &shortener{
		demo.Shortener{},
		base_,
	}, nil
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

	body := struct {
		Url string `json:"url"`
	}{}
	body.Url = s.base.JoinPath(slug).String()

	if err := json.NewEncoder(w).Encode(body); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

type handler func(w http.ResponseWriter, r *http.Request) (int, error)

func handle(f handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if code, err := f(w, r); err != nil {
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
