package runMocks

import (
	"log"
	"strings"
	"sync"

	"github.com/elijahthis/ngatex/internal/mockservice"
	"github.com/elijahthis/ngatex/internal/mockservice2"
	"github.com/elijahthis/ngatex/pkg/config"
)

func Run() {
	configData, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Unable to load config file: %v", err)
	}

	var wg sync.WaitGroup

	for _, upstream := range configData.Services.ServiceA.Upstreams {
		splitPort := strings.Split(upstream, ":")
		port := splitPort[len(splitPort)-1]

		wg.Add(1)
		go mockservice.RunMultiple(port, &wg)
	}

	for _, upstream := range configData.Services.ServiceB.Upstreams {
		splitPort := strings.Split(upstream, ":")
		port := splitPort[len(splitPort)-1]

		wg.Add(1)
		go mockservice2.RunMultiple(port, &wg)
	}

	wg.Wait()
}
