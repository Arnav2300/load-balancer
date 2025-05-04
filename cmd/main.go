package main

import (
	"log"
	"net/http"
	"time"

	"load-balancer/internal/config"
	"load-balancer/internal/handler"
	"load-balancer/internal/loadbalancer"
	"load-balancer/internal/logger"

	"go.uber.org/zap"
)

func main() {
	yamlBytes, _ := config.LoadConfig()
	cfg, _ := config.ParseConfig(yamlBytes)

	zapLogger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	lb := loadbalancer.NewLoadBalancer()
	gateway := handler.NewGateway(cfg, lb)

	loggedHandler := logger.LoggingMiddleware(zapLogger)(gateway)

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      loggedHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	zapLogger.Info("starting load balancer", zap.String("port", cfg.Port))
	if err := srv.ListenAndServe(); err != nil {
		zapLogger.Fatal("failed to start load balancer", zap.Error(err))
	}
}
