package main

import (
	"log"
	"net/http"
	"time"

	"load-balancer/internal/config"
	"load-balancer/internal/handler"
	"load-balancer/internal/healthcheck"
	"load-balancer/internal/loadbalancer"
	"load-balancer/internal/logger"

	"go.uber.org/zap"
)

func main() {
	yamlBytes, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg, err := config.ParseConfig(yamlBytes)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	zapLogger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	logger.Inject(zapLogger)
	hc := healthcheck.NewHealthChecker(5 * time.Second)
	lb := loadbalancer.NewLoadBalancer()
	lb.SetHealthChecker(hc)
	gateway := handler.NewGateway(cfg, lb)
	hc.Start(cfg.Routes)
	loggedHandler := logger.LoggingMiddleware(zapLogger)(gateway)

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      loggedHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	zapLogger.Info("starting load balancer", zap.String("port", cfg.Port), zap.String("env", cfg.Env))
	if err := srv.ListenAndServe(); err != nil {
		zapLogger.Fatal("failed to start load balancer", zap.Error(err))
	}
}
