package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/elijahthis/ngatex/pkg/config"
	"github.com/elijahthis/ngatex/pkg/health"
	"github.com/elijahthis/ngatex/pkg/loadbalancer"
	"github.com/elijahthis/ngatex/pkg/router"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	log.Printf("Loading config from %s", *configPath)

	configData, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Gateway Main: %v", err)
	}

	r := router.New()
	routeMap := config.BuildRouteServiceMap(configData)

	for _, route := range configData.Routes {
		route := route
		service := routeMap[route.Path]
		var lb loadbalancer.Balancer

		switch service.LoadBalancingPolicy {
		case "round-robin":
			lb, err = loadbalancer.NewRoundRobin(service.Upstreams)
		case "weighted-round-robin":
			lb, err = loadbalancer.NewWeightedRoundRobin(service.Upstreams)
		case "least-connections":
			lb, err = loadbalancer.NewLeastConnections(service.Upstreams)
		}

		if err != nil {
			log.Fatalf("failed to initialize load balancer for %s: %v", route.Path, err)
		}

		health.StartActiveServiceChecks(lb.GetUpstreams(), 10*time.Second)

		r.AddRoute(route.Path, service, lb)
	}

	log.Println("Gateway running on :8080")
	http.ListenAndServe(":8080", r.Router)

}
