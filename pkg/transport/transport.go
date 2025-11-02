package transport

import (
	"net"
	"net/http"
	"time"
)

type TransportConfig struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	DialTimeout         time.Duration
	TLSHandshakeTimeout time.Duration
	// RequestTimeout      time.Duration
	// MaxRetries          int
}

func NewGatewayTransport(cfg TransportConfig) *http.Transport {
	return &http.Transport{
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:     cfg.IdleConnTimeout,
		TLSHandshakeTimeout: cfg.TLSHandshakeTimeout,
		DialContext: (&net.Dialer{
			Timeout:   cfg.DialTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		// implement later, will try Roundtripper
		// RequestTimeout:      cfg.RequestTimeout,
		// MaxRetries:          cfg.MaxIdleConns,
	}

}
