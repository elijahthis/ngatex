package loadbalancer

import (
	"errors"
	"log"
	"net/url"
	"sync/atomic"
)

type RoundRobin struct {
	upstreams []*Upstream
	next      atomic.Uint64
}

func NewRoundRobin(upStreamStrings []string) (*RoundRobin, error) {
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

	return &RoundRobin{
		upstreams: upstreams,
	}, nil
}

func (rr *RoundRobin) GetUpstreams() []*Upstream {
	return rr.upstreams
}

func (rr *RoundRobin) Next() (*Upstream, error) {
	numOfUpstreams := uint64(len(rr.upstreams))

	for range numOfUpstreams {
		n := rr.next.Add(1)
		idx := (n - 1) % numOfUpstreams

		upstream := rr.upstreams[idx]
		if upstream.IsAlive() {
			log.Printf("RR chose %s", upstream.URL.String())
			return upstream, nil
		}
	}

	return nil, errors.New("no healthy upstreams available")

}
