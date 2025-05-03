package service

import (
	"errors"
	"lb/internal/models"
	"lb/internal/models/dto"
	"lb/pkg/healthcheck"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type LoadBalancerService struct {
	serverPool *models.ServerPool
	health     *healthcheck.HealthChecker
}

func NewLoadBalancerService() *LoadBalancerService {
	return &LoadBalancerService{
		serverPool: &models.ServerPool{},
		health:     healthcheck.NewHealthChecker(),
	}
}

func (lb *LoadBalancerService) AddBackend(
	backendDTO *dto.AddBackendRequest,
	proxyErrorHandler func(w http.ResponseWriter, r *http.Request, err error),
) error {
	u, err := url.Parse(backendDTO.ServerURL)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.ErrorHandler = proxyErrorHandler
	backend := models.NewBackend(u, proxy)

	lb.serverPool.AddBackend(backend)
	return nil
}

func (lb *LoadBalancerService) DeleteBackend(backendDTO *dto.DeleteBackendRequest) error {
	u, err := url.Parse(backendDTO.ServerURL)
	if err != nil {
		return err
	}

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

func (lb *LoadBalancerService) ServerPool() *models.ServerPool {
	return lb.serverPool
}

func (lb *LoadBalancerService) HealthChecker() *healthcheck.HealthChecker {
	return lb.health
}
