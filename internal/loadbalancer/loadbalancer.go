package loadbalancer

import (
	"load-balancer/internal/healthcheck"
	"sync"
)

type LoadBalancer struct {
	mu            sync.Mutex
	counters      map[string]int
	healthChecker *healthcheck.HealthChecker
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		counters: make(map[string]int),
	}
}

func (lb *LoadBalancer) SetHealthChecker(hc *healthcheck.HealthChecker) {
	lb.healthChecker = hc
}

func (lb *LoadBalancer) SelectBackend(pathPrefix string, backends []string) (string, bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(backends) == 0 {
		return "", false
	}
	healthy := make([]string, 0, len(backends))
	for _, b := range backends {
		if lb.healthChecker == nil || lb.healthChecker.IsHealthy(b) {
			healthy = append(healthy, b)
		}
	}
	if len(healthy) == 0 {
		return "", false
	}
	idx := lb.counters[pathPrefix] % len(healthy)
	backend := healthy[idx]
	lb.counters[pathPrefix] = (lb.counters[pathPrefix] + 1) % len(healthy)

	return backend, true
}
