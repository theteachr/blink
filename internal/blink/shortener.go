package blink

import "net/url"

type Shortener interface {
	Shorten(*url.URL) (string, error)
	Resolve(string) (*url.URL, error)
}
