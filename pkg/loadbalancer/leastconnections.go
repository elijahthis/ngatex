package loadbalancer

import (
	"errors"
	"net/url"
	"sync"
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
	}, nil
}

func (lc *LeastConnections) Next() (*Upstream, error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// numOfUpstreams := uint64(len(rr.upstreams))
	var best *Upstream
	minConn := int64(-1)

	for _, u := range lc.upstreams {
		if !u.IsAlive() {
			continue
		}

		if best == nil || u.GetActiveConn() < minConn {
			best = u
		}
	}

	if best == nil {
		return nil, errors.New("no healthy upstreams available")
	}

	return best, nil
}
