package models

import (
	"log"
	"net/url"
	"sync/atomic"
)

type ServerPool struct {
	Backends []*Backend
	Current  uint64
}

func (s *ServerPool) NextIndex() int {
	if len(s.Backends) == 0 {
		return -1
	}
	return int(atomic.AddUint64(&s.Current, uint64(1))) % len(s.Backends)
}

func (s *ServerPool) AddBackend(backend *Backend) {
	s.Backends = append(s.Backends, backend)
}

func (s *ServerPool) MarkBackendStatus(serverUrl *url.URL, alive bool) {
	for _, b := range s.Backends {
		if b.URL.String() == serverUrl.String() {
			b.SetAlive(alive)
			return
		}
	}
}

func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	if next == -1 {
		log.Println("No backends available in the server pool")
		return nil
	}

	cycleLen := len(s.Backends) + next
	for i := next; i < cycleLen; i++ {
		idx := i % len(s.Backends)
		if s.Backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.Current, uint64(idx))
			}
			return s.Backends[idx]
		}
	}
	return nil
}
