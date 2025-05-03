package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func startServer(wg *sync.WaitGroup, port int, serverName string) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello from %s!", serverName)
		if err != nil {
			return
		}
	})

	log.Printf("Starting server %s on port %d\n", serverName, port)

	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server %s failed: %v\n", serverName, err)
	}
}

func main() {
	servers := []struct {
		port       int
		serverName string
	}{
		{8081, "Server 1"},
		{8082, "Server 2"},
		{8083, "Server 3"},
	}

	var wg sync.WaitGroup

	for _, srv := range servers {
		wg.Add(1)
		go startServer(&wg, srv.port, srv.serverName)
	}

	wg.Wait()
}
