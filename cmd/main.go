package main

import (
	"fmt"

	"load-balancer/internal/config"
	"load-balancer/internal/proxy"
)

func main() {
	yamlBytes, _ := config.LoadConfig()
	cfg, _ := config.ParseConfig(yamlBytes)
	fmt.Println(cfg.Routes[2])
	proxy.Server()
}
