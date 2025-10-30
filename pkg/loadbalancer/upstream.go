package loadbalancer

import (
	"net/url"
	"sync/atomic"
)

type Upstream struct {
	URL        *url.URL
	Weight     int32
	activeConn atomic.Int64
	healthy    atomic.Bool
}

func (u *Upstream) IsAlive() bool {
	return u.healthy.Load()
}

func (u *Upstream) SetAlive(alive bool) {
	u.healthy.Store(alive)
}

func (u *Upstream) IncActiveConn() {
	u.activeConn.Add(1)
}

func (u *Upstream) DecActiveConn() {
	u.activeConn.Add(-1)
}

func (u *Upstream) GetActiveConn() int64 {
	return u.activeConn.Load()
}
