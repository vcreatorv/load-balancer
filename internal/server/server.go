package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
		},
	}
}

func (s *Server) SetupRoutes(routeConfig func(*http.ServeMux)) {
	subrouter := http.NewServeMux()
	mainRouter := http.NewServeMux()
	mainRouter.Handle("/loadbalancer/", http.StripPrefix("/loadbalancer", subrouter))

	routeConfig(subrouter)

	s.httpServer.Handler = mainRouter
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
