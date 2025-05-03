package healthcheck

import (
	"load-balancer/internal/config"
	"net/http"
	"sync"
	"time"
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
						hc.checkBackend(b)
					}(backend)
				}
			}
			wg.Wait()
			time.Sleep(hc.interval)
		}
	}()
}

func (hc *HealthChecker) checkBackend(backend string) {
	hc.mu.RLock()
	status, exists := hc.healthStatus[backend]
	hc.mu.RUnlock()

	if !exists {
		status = &backendStatus{}
	}

	if time.Now().Before(status.NextRetryAfter) {
		return
	}

	resp, err := hc.client.Get(backend + "/health")
	healthy := err == nil && resp.StatusCode == http.StatusOK

	hc.mu.Lock()
	defer hc.mu.Unlock()
	// hc.healthStatus[backend] = status

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
	}

	hc.healthStatus[backend] = status
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
func (hc *HealthChecker) IsHealthy(backend string) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return hc.healthStatus[backend].Healthy
}
