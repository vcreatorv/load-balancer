package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
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
	Connections  uint64
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

func (b *Backend) IncrementConnections() {
	atomic.AddUint64(&b.Connections, uint64(1))
}

func (b *Backend) DecrementConnections() {
	atomic.AddUint64(&b.Connections, ^uint64(0))
}

func (b *Backend) GetConnections() uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Connections
}

func NewBackend(serverURL *url.URL, proxy *httputil.ReverseProxy) *Backend {
	return &Backend{
		URL:          serverURL,
		Alive:        true,
		ReverseProxy: proxy,
		Connections:  0,
	}
}
