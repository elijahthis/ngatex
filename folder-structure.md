api-gateway/
├── cmd/
│   ├── gateway/
│   │   └── main.go                # Entry point for the main gateway binary
│   └── admin/
│       └── main.go                # Optional separate admin server binary
│
├── internal/
│   ├── config/
│   │   └── config.go              # Loads YAML/JSON config into structs
│   │
│   ├── router/
│   │   ├── router.go              # Route registry and matching logic
│   │   └── trie.go                # Optional high-performance router implementation
│   │
│   ├── proxy/
│   │   └── proxy.go               # Core reverse proxy logic (load balancing, retries)
│   │
│   ├── middleware/
│   │   ├── auth.go                # JWT / API key authentication
│   │   ├── ratelimit.go           # Token bucket / leaky bucket limiter
│   │   ├── cache.go               # In-memory or Redis cache
│   │   ├── circuitbreaker.go      # Simple fail-fast logic
│   │   └── chain.go               # Middleware chaining logic
│   │
│   ├── observability/
│   │   ├── metrics.go             # Prometheus metrics
│   │   ├── tracing.go             # OpenTelemetry integration
│   │   └── logging.go             # Zap/Zerolog setup
│   │
│   ├── admin/
│   │   └── admin_server.go        # Admin API: /admin/routes, /admin/metrics, etc.
│   │
│   ├── store/
│   │   ├── store.go               # BoltDB or etcd abstraction layer
│   │   └── models.go              # Route, Upstream, and Service models
│   │
│   └── loadbalancer/
│       └── balancer.go            # RoundRobin, LeastConn, etc.
│
├── pkg/
│   └── utils/
│       └── helpers.go             # Shared helper functions
│
├── configs/
│   └── config.yaml                # Example gateway configuration
│
├── tests/
│   └── integration_test.go        # Integration and load tests
│
├── go.mod
├── go.sum
└── README.md
