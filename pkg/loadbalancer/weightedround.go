package loadbalancer

import (
	"errors"
	"net/url"
	"sync"

	"github.com/rs/zerolog/log"
)

type WeightedRoundRobin struct {
	upstreams []*Upstream
	mu        *sync.Mutex
}

func NewWeightedRoundRobin(upStreamStrings []string) (*WeightedRoundRobin, error) {
	upstreams := make([]*Upstream, len(upStreamStrings))
	for i, u := range upStreamStrings {
		parsedURL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		upstreams[i] = &Upstream{
			URL:           parsedURL,
			Weight:        int32(i%2) + 1,
			currentWeight: 0,
		}
		upstreams[i].SetAlive(true)
	}

	return &WeightedRoundRobin{
		upstreams: upstreams,
		mu:        &sync.Mutex{},
	}, nil
}

func (wrr *WeightedRoundRobin) GetUpstreams() []*Upstream {
	return wrr.upstreams
}

func (wrr *WeightedRoundRobin) Next() (*Upstream, error) {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	totalWeights := int64(0)

	var selectedUpstream *Upstream

	for _, u := range wrr.upstreams {
		if !u.IsAlive() {
			continue
		}
		u.currentWeight += u.Weight
		totalWeights += int64(u.Weight)

		if selectedUpstream == nil || u.currentWeight > selectedUpstream.currentWeight {
			selectedUpstream = u
		}
	}

	if selectedUpstream == nil {
		return nil, errors.New("no healthy upstreams available")
	}

	selectedUpstream.currentWeight -= int32(totalWeights)
	log.Info().Msgf("WRR chose %s", selectedUpstream.URL.String())

	return selectedUpstream, nil
}
