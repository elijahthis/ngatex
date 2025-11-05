package loadbalancer

import (
	"errors"
	"net/url"
	"sync"

	"github.com/rs/zerolog/log"
)

type LeastConnections struct {
	upstreams []*Upstream
	mu        *sync.Mutex
}

func NewLeastConnections(upStreamStrings []string) (*LeastConnections, error) {
	upstreams := make([]*Upstream, len(upStreamStrings))
	for i, u := range upStreamStrings {
		parsedURL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		upstreams[i] = &Upstream{
			URL:    parsedURL,
			Weight: 1,
		}
		upstreams[i].SetAlive(true)
	}

	return &LeastConnections{
		upstreams: upstreams,
		mu:        &sync.Mutex{},
	}, nil
}

func (lc *LeastConnections) GetUpstreams() []*Upstream {
	return lc.upstreams
}

func (lc *LeastConnections) Next() (*Upstream, error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	var best *Upstream
	minConn := int64(1<<63 - 1)

	for _, u := range lc.upstreams {
		if !u.IsAlive() {
			continue
		}

		if best == nil || u.GetActiveConn() < minConn {
			best = u
			minConn = u.GetActiveConn()
		}
	}

	if best == nil {
		return nil, errors.New("no healthy upstreams available")
	}

	log.Info().Msgf("LC chose %s", best.URL.String())

	return best, nil
}
