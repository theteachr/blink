package demo

import "net/url"

type Shortener struct{}

func (Shortener) Shorten(*url.URL) (string, error) {
	return "demo", nil
}

func (Shortener) Resolve(slug string) (*url.URL, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "theteachr.github.io",
	}

	return u, nil
}
