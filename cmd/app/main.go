package main

import (
	"errors"
	handler "lb/internal/delivery/http"
	"lb/internal/server"
	"lb/internal/usecase/service"
	"lb/pkg/healthcheck"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const addr = "localhost:8080"

func main() {
	loadBalancerUC := service.NewLoadBalancerService()
	loadBalancerHandler := handler.NewLoadBalancerHandler(loadBalancerUC)

	healthcheck.StartHealthCheck(2*time.Minute, loadBalancerUC.ServerPool(), loadBalancerUC.HealthChecker())

	srv := server.NewServer(addr)
	srv.SetupRoutes(func(r *http.ServeMux) {
		loadBalancerHandler.Configure(r)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		log.Printf("Завершение работы сервера приложения... : %v", sig)
		if err := srv.Stop(); err != nil {
			log.Fatalf("Не удалось остановить сервер приложения: %v", err)
		}
	}()

	log.Printf("Запуск сервера по адресу %s", addr)
	if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
