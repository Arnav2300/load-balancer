package main

import (
	"fmt"
	"log"
	"net/http"

	"load-balancer/internal/config"
	"load-balancer/internal/handler"
	"load-balancer/internal/loadbalancer"
)

func main() {
	yamlBytes, _ := config.LoadConfig()
	cfg, _ := config.ParseConfig(yamlBytes)
	fmt.Println(cfg.Port)
	lb := loadbalancer.NewLoadBalancer()
	gateway := handler.NewGateway(cfg, lb)
	if err := http.ListenAndServe(cfg.Port, gateway); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
