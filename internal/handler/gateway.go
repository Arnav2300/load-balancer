package handler

import (
	"load-balancer/internal/config"
	"load-balancer/internal/loadbalancer"
	"load-balancer/internal/proxy"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"
	"sync"
)

type Gateway struct {
	config       *config.Config
	loadBalancer *loadbalancer.LoadBalancer
	proxies      map[string]*httputil.ReverseProxy
	proxiesMU    sync.RWMutex
}

func NewGateway(cfg *config.Config, lb *loadbalancer.LoadBalancer) *Gateway {
	return &Gateway{
		config:       cfg,
		loadBalancer: lb,
		proxies:      make(map[string]*httputil.ReverseProxy),
	}
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := g.matchRoute(r.URL.Path)
	if route == nil {
		http.NotFound(w, r)
		return
	}
	if !slices.Contains(route.Methods, r.Method) && len(route.Methods) > 0 {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	backendURL, ok := g.loadBalancer.SelectBackend(route.PathPrefix, route.Backends)
	if !ok {
		http.Error(w, "No backends available", http.StatusBadGateway)
		return
	}

	target, err := url.Parse(backendURL)
	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}
	reverseProxy := g.getOrCreateProxy(target)
	handler := proxy.ProxyRequestHandler(reverseProxy, target, route.PathPrefix)
	handler(w, r)

}

func (g *Gateway) matchRoute(path string) *config.Route {
	for _, route := range g.config.Routes {
		if strings.HasPrefix(path, route.PathPrefix) {
			return &route
		}
	}
	return nil
}
func (g *Gateway) getOrCreateProxy(target *url.URL) *httputil.ReverseProxy {
	g.proxiesMU.RLock()
	existingProxy, exists := g.proxies[target.String()]
	g.proxiesMU.RUnlock()
	if exists {
		return existingProxy
	}
	newProxy := proxy.NewReverseProxy(target)
	g.proxiesMU.Lock()
	g.proxies[target.String()] = newProxy
	g.proxiesMU.Unlock()

	return newProxy
}
