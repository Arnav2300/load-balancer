package proxy

import (
	"fmt"
	"io"
	"load-balancer/internal/config"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var backendIndex = make(map[string]int)
var mu sync.Mutex

func Server() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Host)
		fmt.Println(req.Method)
		fmt.Println(req.URL)
		str := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
		var pathPrefix string
		if len(str) > 0 {
			pathPrefix = str[0]
			fmt.Println(pathPrefix)
		} else {
			fmt.Println("no path prefix found")
		}
		yamlBytes, _ := config.LoadConfig()
		cfg, _ := config.ParseConfig(yamlBytes)
		fmt.Println(cfg.Routes)
		for _, route := range cfg.Routes {
			if route.PathPrefix == "/"+pathPrefix {
				mu.Lock()
				index := backendIndex[route.PathPrefix]
				url, _ := url.Parse(route.Backends[index])
				backendIndex[route.PathPrefix] = (index + 1) % len(route.Backends)
				mu.Unlock()
				fmt.Println(backendIndex[route.PathPrefix])
				proxy := NewProxy(url)
				handler := ProxyRequestHandler(proxy, url, route.PathPrefix)
				handler(w, req)
				return
			}
		}
		io.WriteString(w, "hello world\n")
	}

	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
