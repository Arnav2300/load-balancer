package proxy

import (
	"fmt"
	"load-balancer/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	target.Path = ""
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req) // Applies the default modifications
		fullURL := fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path)
		if req.URL.RawQuery != "" {
			fullURL += "?" + req.URL.RawQuery
		}
		logger.Get().Info("forwarding to", zap.String("url", fullURL))
	}
	return proxy
}
func ProxyRequestHandler(proxy *httputil.ReverseProxy, url *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// fmt.Println(r.Host, r.URL.Host, url.Host)
		r.URL.Host = url.Host
		// fmt.Println("host->", url.Host)
		r.URL.Scheme = url.Scheme
		// fmt.Println("Scheme->", url.Scheme)
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Host = url.Host
		// path := r.URL.Path
		// r.URL.Path = "/" + strings.TrimPrefix(strings.TrimPrefix(path, endpoint), "/")

		proxy.ServeHTTP(w, r)

	}
}
