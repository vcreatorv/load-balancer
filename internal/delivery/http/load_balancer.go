package http

import (
	"context"
	"encoding/json"
	"lb/internal/delivery/http/utils"
	"lb/internal/models"
	"lb/internal/models/dto"
	"lb/internal/usecase"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	MAX_ATTEMPTS = 3
	MAX_RETRIES  = 3
)

type LoadBalancerHandler struct {
	LoadBalancerUC usecase.LoadBalancer
}

func (h *LoadBalancerHandler) Configure(r *http.ServeMux) {
	loadBalancerMux := http.NewServeMux()

	loadBalancerMux.HandleFunc("POST /add", h.AddBackend)
	loadBalancerMux.HandleFunc("POST /delete", h.DeleteBackend)

	r.Handle("/backend/", http.StripPrefix("/backend", loadBalancerMux))
	r.HandleFunc("/", h.ForwardRequest)
}

func NewLoadBalancerHandler(loadBalancerUC usecase.LoadBalancer) *LoadBalancerHandler {
	return &LoadBalancerHandler{
		LoadBalancerUC: loadBalancerUC,
	}
}

func (h *LoadBalancerHandler) AddBackend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, models.NewError(models.ErrForbidden, "method is not allowed"))
		return
	}
	var backendDTO dto.AddBackendRequest
	if err := json.NewDecoder(r.Body).Decode(&backendDTO); err != nil {
		utils.WriteError(w, models.NewError(models.ErrBadRequest, err.Error()))
		return
	}

	if err := backendDTO.Validate(); err != nil {
		utils.WriteError(w, models.NewError(models.ErrBadRequest, err.Error()))
		return
	}

	u, err := url.Parse(backendDTO.ServerURL)
	if err != nil {
		utils.WriteError(w, models.NewError(models.ErrBadRequest, err.Error()))
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = h.ProxyErrorHandler
	if err := h.LoadBalancerUC.AddBackend(u, proxy); err != nil {
		utils.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LoadBalancerHandler) DeleteBackend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, models.NewError(models.ErrForbidden, "method is not allowed"))
		return
	}

	var backendDTO dto.DeleteBackendRequest
	if err := json.NewDecoder(r.Body).Decode(&backendDTO); err != nil {
		utils.WriteError(w, models.NewError(models.ErrBadRequest, err.Error()))
		return
	}

	if err := backendDTO.Validate(); err != nil {
		utils.WriteError(w, models.NewError(models.ErrBadRequest, err.Error()))
		return
	}

	if err := h.LoadBalancerUC.DeleteBackend(&backendDTO); err != nil {
		utils.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LoadBalancerHandler) ForwardRequest(w http.ResponseWriter, r *http.Request) {
	attempts := h.getAttemptsFromContext(r)
	if attempts > MAX_ATTEMPTS {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		utils.WriteError(w, models.NewError(models.ErrInternal, "service not available"))
		return
	}

	pool := h.LoadBalancerUC.ServerPool()
	backend := pool.GetNextPeer()
	if backend != nil {
		backend.ReverseProxy.ServeHTTP(w, r)
		return
	}
	utils.WriteError(w, models.NewError(models.ErrInternal, "service not available"))
}

func (h *LoadBalancerHandler) getAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(models.Attempts).(int); ok {
		return attempts
	}
	return 0
}

func (h *LoadBalancerHandler) getRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(models.Retry).(int); ok {
		return retry
	}
	return 0
}

func (h *LoadBalancerHandler) ProxyErrorHandler(w http.ResponseWriter, r *http.Request, e error) {
	log.Printf("PROXY ERROR: [%s] %s\n", r.URL.Host, e.Error())

	retries := h.getRetryFromContext(r)
	if retries < MAX_RETRIES {
		select {
		case <-time.After(10 * time.Millisecond):
			ctx := context.WithValue(r.Context(), models.Retry, retries+1)
			r = r.WithContext(ctx)
			h.ForwardRequest(w, r)
		}
		return
	}

	h.LoadBalancerUC.MarkBackendStatus(r.URL, false)

	attempts := h.getAttemptsFromContext(r)
	ctx := context.WithValue(r.Context(), models.Attempts, attempts+1)
	h.ForwardRequest(w, r.WithContext(ctx))
}
