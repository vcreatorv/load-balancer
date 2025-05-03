package usecase

import (
	"lb/internal/models"
	"lb/internal/models/dto"
	"lb/pkg/healthcheck"
	"net/http"
	"net/url"
)

type LoadBalancer interface {
	AddBackend(backendDTO *dto.AddBackendRequest, proxyErrorHandler func(w http.ResponseWriter, r *http.Request, err error)) error
	DeleteBackend(backendDTO *dto.DeleteBackendRequest) error
	ServerPool() *models.ServerPool
	HealthChecker() *healthcheck.HealthChecker
	MarkBackendStatus(serverUrl *url.URL, alive bool)
}
