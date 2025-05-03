package main

import (
	"errors"
	"flag"
	"fmt"
	handler "lb/internal/delivery/http"
	"lb/internal/server"
	"lb/internal/usecase/service"
	"lb/pkg/healthcheck"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	port    = flag.Int("port", 8080, "listen address")
	servers = &AddrList{}
)

type AddrList []string

func (a *AddrList) String() string {
	return strings.Join(*a, ",")
}

func (a *AddrList) Set(in string) error {
	for _, addr := range strings.Split(in, ",") {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}

		if _, err := url.Parse(addr); err != nil {
			return fmt.Errorf("incorrect URL %q: %v", addr, err)
		}
		*a = append(*a, addr)
	}
	return nil
}

func init() {
	flag.Var(servers, "servers", "list of servers, ex. http://127.0.0.1:8081,http://127.0.0.1:8082")
	flag.Parse()
}

// go run .\cmd\app\main.go --port=8090 --servers="http://127.0.0.1:8081,http://127.0.0.1:8082"

func main() {
	loadBalancerUC := service.NewLoadBalancerService()
	loadBalancerHandler := handler.NewLoadBalancerHandler(loadBalancerUC)

	for _, addr := range *servers {
		u, err := url.Parse(addr)
		if err != nil {
			log.Printf("Error parsing URL %s: %v", addr, err)
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ErrorHandler = loadBalancerHandler.ProxyErrorHandler

		if err := loadBalancerUC.AddBackend(u, proxy); err != nil {
			log.Printf("Error adding backend: %v", err)
		} else {
			log.Printf("Added backend: %v", u)
		}
	}

	healthcheck.StartHealthCheck(2*time.Minute, loadBalancerUC.ServerPool(), loadBalancerUC.HealthChecker())

	listenPort := fmt.Sprintf(":%d", *port)
	srv := server.NewServer(listenPort)
	srv.SetupRoutes(func(r *http.ServeMux) {
		loadBalancerHandler.Configure(r)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		log.Printf("Shutting down the application server... : %v", sig)
		if err := srv.Stop(); err != nil {
			log.Fatalf("Failed to stop the application server: %v", err)
		}
	}()

	log.Printf("Starting server at port %s", listenPort)
	if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
