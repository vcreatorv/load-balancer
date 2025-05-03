package usecase

import (
	"lb/internal/models"
	"lb/internal/models/dto"
	"lb/pkg/healthcheck"
	"net/http/httputil"
	"net/url"
)

type LoadBalancer interface {
	AddBackend(url *url.URL, proxy *httputil.ReverseProxy) error
	DeleteBackend(backendDTO *dto.DeleteBackendRequest) error
	ServerPool() *models.ServerPool
	HealthChecker() *healthcheck.HealthChecker
	MarkBackendStatus(serverUrl *url.URL, alive bool)
}
