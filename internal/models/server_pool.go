package models

import (
	"log"
	"sync/atomic"
)

type BalancingAlgorithm int

const (
	RoundRobin BalancingAlgorithm = iota
	LeastConnections
)

type ServerPool struct {
	Backends  []*Backend
	Current   uint64
	Algorithm BalancingAlgorithm
}

func (s *ServerPool) SetAlgorithm(algorithm BalancingAlgorithm) {
	s.Algorithm = algorithm
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

func (s *ServerPool) GetNextPeer() *Backend {
	switch s.Algorithm {
	case LeastConnections:
		return s.getNextPeerLeastConnections()
	case RoundRobin:
		return s.getNextPeerRoundRobin()
	default:
		return s.getNextPeerRoundRobin()
	}
}

func (s *ServerPool) getNextPeerRoundRobin() *Backend {
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
			s.Backends[idx].IncrementConnections()
			log.Printf(
				"RoundRobin selected backend %s [connections: %d]",
				s.Backends[idx].URL.Host,
				s.Backends[idx].GetConnections(),
			)
			return s.Backends[idx]
		}
	}
	return nil
}

func (s *ServerPool) getNextPeerLeastConnections() *Backend {
	var leastConnBackend *Backend
	minConnections := ^uint64(0) >> 1

	var selectedIdx int
	for i, backend := range s.Backends {
		if backend.IsAlive() {
			conn := backend.GetConnections()
			if conn < minConnections {
				minConnections = conn
				leastConnBackend = backend
				selectedIdx = i
			}
		}
	}

	if leastConnBackend != nil {
		atomic.StoreUint64(&s.Current, uint64(selectedIdx))
		leastConnBackend.IncrementConnections()
		log.Printf(
			"LeastConnections selected backend %s [connections: %d]",
			leastConnBackend.URL.Host,
			leastConnBackend.GetConnections(),
		)
		return leastConnBackend
	}

	log.Println("No backends available in the server pool")
	return nil
}
