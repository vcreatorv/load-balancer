package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

const (
	Attempts int = iota
	Retry
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mu           sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	b.Alive = alive
	b.mu.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Alive
}

func NewBackend(serverURL *url.URL, proxy *httputil.ReverseProxy) *Backend {
	return &Backend{
		URL:          serverURL,
		Alive:        true,
		ReverseProxy: proxy,
	}
}
