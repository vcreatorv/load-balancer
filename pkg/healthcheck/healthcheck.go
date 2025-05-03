package healthcheck

import (
	"lb/internal/models"
	"log"
	"net"
	"net/url"
	"time"
)

type HealthChecker struct{}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

func (h *HealthChecker) Check(pool *models.ServerPool) {
	for _, b := range pool.Backends {
		status := "UP"
		alive := h.IsBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "DOWN"
		}
		log.Printf("Backend %s [status=%s]\n", b.URL, status)
	}
}

func (h *HealthChecker) IsBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	defer conn.Close()
	return true
}
