package healthcheck

import (
	"fmt"
	"load-balancer/internal/config"
	"load-balancer/internal/logger"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type backendStatus struct {
	Healthy        bool
	Failures       int
	NextRetryAfter time.Time
}

type HealthChecker struct {
	healthStatus map[string]*backendStatus
	mu           sync.RWMutex
	interval     time.Duration
	client       *http.Client
}

func NewHealthChecker(interval time.Duration) *HealthChecker {
	return &HealthChecker{
		healthStatus: make(map[string]*backendStatus),
		interval:     interval,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (hc *HealthChecker) Start(routes []config.Route) {
	go func() {
		for {
			var wg sync.WaitGroup
			for _, route := range routes {
				for _, backend := range route.Backends {
					wg.Add(1)
					go func(b string) {
						defer wg.Done()
						hc.checkBackend(b + route.PathPrefix + route.Health)
					}(backend)
				}
			}
			wg.Wait()
			time.Sleep(hc.interval)
		}
	}()
}

func (hc *HealthChecker) checkBackend(backend string) {
	hc.mu.Lock()
	status, exists := hc.healthStatus[backend]
	if !exists {
		status = &backendStatus{}
	}
	hc.mu.Unlock()

	if time.Now().Before(status.NextRetryAfter) {
		return
	}

	resp, err := hc.client.Get(backend)
	healthy := err == nil && resp.StatusCode == http.StatusOK

	hc.mu.Lock()
	defer hc.mu.Unlock()
	// hc.healthStatus[backend] = status
	prevHealthy := status.Healthy

	if healthy {
		status.Healthy = true
		status.Failures = 0
		status.NextRetryAfter = time.Time{}
	} else {
		status.Healthy = false
		status.Failures++
		backoffDuration := time.Duration(1<<status.Failures) * time.Second
		if backoffDuration > 5*time.Minute {
			backoffDuration = 5 * time.Minute
		}
		status.NextRetryAfter = time.Now().Add(backoffDuration)
		// logger.Get().Info("health check", zap.String("DEAD", backend))
	}
	if prevHealthy != status.Healthy {
		if status.Healthy {
			logger.Get().Info("health check", zap.String("ALIVE", backend))
		} else {
			logger.Get().Warn("health check", zap.String("DEAD", backend))
		}
	}
	hc.healthStatus[backend] = status
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
func (hc *HealthChecker) IsHealthy(backend string) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	status, ok := hc.healthStatus[backend]
	for key, value := range hc.healthStatus {
		fmt.Println("Key:", key, "Value:", value)
	}
	// fmt.Println(backend, " ", status.Healthy)
	return ok && status.Healthy
}
