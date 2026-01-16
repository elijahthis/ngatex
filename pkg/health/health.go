package health

import (
	"context"
	"net/http"
	"time"

	"github.com/elijahthis/ngatex/pkg/loadbalancer"
	"github.com/rs/zerolog/log"
)

var healthClient = http.Client{
	Timeout: 4 * time.Second,
}

func StartActiveServiceChecks(upstreams []*loadbalancer.Upstream, sleepDuration time.Duration, ctx context.Context) {
	for _, upstream := range upstreams {
		go func() {
			ticker := time.NewTicker(sleepDuration)
			defer ticker.Stop()

			upstream := upstream
			for {
				select {
				case <-ctx.Done():
					log.Info().Str("url", upstream.URL.String()).Msg("Stopping health check (context cancelled)")
					return
				case <-ticker.C:
					resp, err := healthClient.Get(upstream.URL.String() + "/health")
					alive := err == nil && resp.StatusCode == http.StatusOK
					upstream.SetAlive(alive)
					if resp != nil {
						resp.Body.Close()
					}
				}

			}

		}()
	}

}
