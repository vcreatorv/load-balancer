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
	Mutex        sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (b *Backend) SetAlive(alive bool) {
	b.Mutex.Lock()
	b.Alive = alive
	b.Mutex.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	return b.Alive
}

func NewBackend(serverURL *url.URL, proxy *httputil.ReverseProxy) *Backend {
	return &Backend{
		URL:          serverURL,
		Alive:        true,
		ReverseProxy: proxy,
	}
}
