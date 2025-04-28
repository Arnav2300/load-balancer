package handler

import (
	"load-balancer/internal/config"
	"load-balancer/internal/loadbalancer"
	"load-balancer/internal/proxy"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

type Gateway struct {
	config       *config.Config
	loadBalancer *loadbalancer.LoadBalancer
}

func NewGateway(cfg *config.Config, lb *loadbalancer.LoadBalancer) *Gateway {
	return &Gateway{
		config:       cfg,
		loadBalancer: lb,
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
	reverseProxy := proxy.NewReverseProxy(target)
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
