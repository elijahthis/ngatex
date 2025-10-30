package main

import (
	"net/http"

	"github.com/elijahthis/ngatex/pkg/router"
)

func main() {
	// Load config file
	// configData, err := config.Load("config.yaml")
	// if err != nil {
	// 	log.Fatalf("Gateway Main: %v", err)
	// }

	// proxyTransport := transport.NewGatewayTransport(transport.TransportConfig{
	// 	MaxIdleConns:        5000,
	// 	MaxIdleConnsPerHost: 1000,
	// 	IdleConnTimeout:     30 * time.Second,
	// 	DialTimeout:         100 * time.Millisecond,
	// 	TLSHandshakeTimeout: 5 * time.Second,
	// })

	r := router.New()

	// Server Multiplexer
	// mux := http.NewServeMux()
	// for _, route := range configData.Routes {
	// 	route := route
	// 	upstreamURL, err := url.Parse(route.Upstream)
	// 	if err != nil {
	// 		log.Fatalf("Gateway Main: %v", err)
	// 	}

	// 	proxy := httputil.NewSingleHostReverseProxy(upstreamURL)
	// 	proxy.Transport = proxyTransport
	// 	handler := http.StripPrefix(route.Path, proxy)
	// 	// proxy.Director = func(req *http.Request) {
	// 	// 	req.URL.Scheme = upstreamURL.Scheme
	// 	// 	req.URL.Host = upstreamURL.Host
	// 	// 	req.Host = upstreamURL.Host
	// 	// 	req.URL.Path = upstreamURL.Path
	// 	// }
	// 	mux.Handle(route.Path, handler)
	// }

	http.ListenAndServe(":8080", r)

}
