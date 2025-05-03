package healthcheck

import (
	"lb/internal/models"
	"log"
	"time"
)

type HealthCheckable interface {
	Check(pool *models.ServerPool)
}

func StartHealthCheck(interval time.Duration, pool *models.ServerPool, checker HealthCheckable) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Starting health check...")
				checker.Check(pool)
				log.Println("Health check completed")
			}
		}
	}()
}
