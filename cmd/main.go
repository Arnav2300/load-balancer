package main

import (
	"fmt"

	"load-balancer/internal/config"
)

func main() {
	yamlBytes, _ := config.LoadConfig()
	cfg, _ := config.ParseConfig(yamlBytes)
	fmt.Println(cfg)
}
