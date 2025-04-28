package loadbalancer

import "sync"

type LoadBalancer struct {
	mu       sync.Mutex
	counters map[string]int
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		counters: make(map[string]int),
	}
}

func (lb *LoadBalancer) SelectBackend(pathPrefix string, backends []string) (string, bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(backends) == 0 {
		return "", false
	}

	idx := lb.counters[pathPrefix] % len(backends)
	backend := backends[idx]
	lb.counters[pathPrefix] = (lb.counters[pathPrefix] + 1) % len(backends)

	return backend, true
}
