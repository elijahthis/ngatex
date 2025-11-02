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
	mockServiceFuncs := []func(string, *sync.WaitGroup){
		mockservice.RunMultiple,
		mockservice2.RunMultiple,
	}

	i := 0
	for _, service := range configData.Services {
		for _, upstream := range service.Upstreams {
			splitPort := strings.Split(upstream, ":")
			port := splitPort[len(splitPort)-1]

			wg.Add(1)
			go mockServiceFuncs[i](port, &wg)
		}
		i++
	}

	wg.Wait()
}
