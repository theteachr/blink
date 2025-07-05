package blink

import "net/url"

type Resolver interface {
	Resolve(string) (*url.URL, error)
}
