package service

import (
	"errors"
	"fmt"
	"lb/internal/models"
	"lb/internal/models/dto"
	"lb/pkg/healthcheck"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type LoadBalancerService struct {
	serverPool *models.ServerPool
	health     *healthcheck.HealthChecker
	mu         sync.RWMutex
	urlMap     map[string]*models.Backend
}

func NewLoadBalancerService() *LoadBalancerService {
	return &LoadBalancerService{
		serverPool: &models.ServerPool{},
		health:     healthcheck.NewHealthChecker(),
		urlMap:     make(map[string]*models.Backend),
	}
}

func (lb *LoadBalancerService) AddBackend(u *url.URL, proxy *httputil.ReverseProxy) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if _, exists := lb.urlMap[u.String()]; exists {
		return fmt.Errorf("backend with url %s already exists", u.String())
	}

	backend := models.NewBackend(u, proxy)
	lb.serverPool.AddBackend(backend)
	lb.urlMap[u.String()] = backend
	return nil
}

func (lb *LoadBalancerService) DeleteBackend(backendDTO *dto.DeleteBackendRequest) error {
	u, err := url.Parse(backendDTO.ServerURL)
	if err != nil {
		return err
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	delete(lb.urlMap, u.String())

	for i, b := range lb.serverPool.Backends {
		if b.URL.String() == u.String() {
			lb.serverPool.Backends = append(lb.serverPool.Backends[:i], lb.serverPool.Backends[i+1:]...)

			if len(lb.serverPool.Backends) > 0 {
				current := atomic.LoadUint64(&lb.serverPool.Current)
				if int(current) >= len(lb.serverPool.Backends) {
					atomic.StoreUint64(&lb.serverPool.Current, 0)
				}
			}
			return nil
		}
	}
	return errors.New("backend not found")
}

func (lb *LoadBalancerService) MarkBackendStatus(serverUrl *url.URL, alive bool) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if backend, exists := lb.urlMap[serverUrl.String()]; exists {
		backend.SetAlive(alive)
	}
}

func (lb *LoadBalancerService) ServerPool() *models.ServerPool {
	return lb.serverPool
}

func (lb *LoadBalancerService) HealthChecker() *healthcheck.HealthChecker {
	return lb.health
}
