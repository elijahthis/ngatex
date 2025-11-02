package health

import (
	"net/http"
	"time"

	"github.com/elijahthis/ngatex/pkg/loadbalancer"
)

var healthClient = http.Client{
	Timeout: 4 * time.Second,
}

func StartActiveServiceChecks(upstreams []*loadbalancer.Upstream, sleepDuration time.Duration) {
	for _, upstream := range upstreams {
		go func() {
			upstream := upstream
			for {
				resp, err := healthClient.Get(upstream.URL.String() + "/health")
				if err != nil {
					upstream.SetAlive(false)
				} else if resp.StatusCode != 200 {
					upstream.SetAlive(false)
				} else {
					upstream.SetAlive(true)
				}

				time.Sleep(sleepDuration)
			}
		}()
	}

}
