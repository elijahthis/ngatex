package loadbalancer

type Balancer interface {
	// Next returns the next available upstream
	Next() (*Upstream, error)
}
