package loadbalancer

type Balancer interface {
	GetUpstreams() []*Upstream
	Next() (*Upstream, error)
}
