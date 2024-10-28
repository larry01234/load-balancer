package main

import (
	"log"
	"net/http"
	"load-balancer/loadbalancer"
	"load-balancer/backend"
)

func main() {
	// Initialize backend nodes with different rate limits
	nodes := []*loadbalancer.Node{
		loadbalancer.NewNode("http://localhost:8081", 1024*1024, 100), // 1MB BPM, 100 RPM
		loadbalancer.NewNode("http://localhost:8082", 512*1024, 50),    // 512KB BPM, 50 RPM
	}

	// Start backend servers
	go backend.StartServer(":8082") // Ensure this is called only once

	// Initialize Load Balancer
	lb := loadbalancer.NewLoadBalancer(nodes)

	// Start Load Balancer server
	log.Println("Load Balancer started at :8080")
	if err := http.ListenAndServe(":8080", lb); err != nil {
		log.Fatalf("Failed to start Load Balancer: %v", err)
	}
}
