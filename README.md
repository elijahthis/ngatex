# Ngatex: High-Performance API Gateway in Go

Ngatex is a lightweight, cloud-native API Gateway built from the ground up in Go. It provides a high-performance entry point for microservices, offering dynamic routing, load balancing, and deep observability.

## 🚀 Key Features

- **High-Performance Proxying:** Built on `net/http` with optimized connection pooling.
- **Dynamic Load Balancing:** Supports Round Robin, Least Connections, and Smooth Weighted Round Robin.
- **Middleware Pipeline:** Pluggable chain for JWT Auth, API Key validation, IP-based Rate Limiting, and Caching.
- **Observability:** Native Prometheus metrics and structured JSON logging with Request ID correlation.
- **Zero-Downtime Reloads:** Support for `SIGHUP` and Admin API reloads using Atomic Pointer Swapping.
- **Health Monitoring:** Active and passive health checks for upstream services.
<!--

## 🏗️ Architecture -->

## 🛠️ Getting Started

### 1. Prerequisites

- Go 1.23+
- Docker (optional)

### 2. Installation

```bash
git clone [https://github.com/your-username/ngatex.git](https://github.com/your-username/ngatex.git)
cd ngatex
make build
```

### 3. Running the System

In one terminal, start the mock backend services:

```bash
make mock
```

In another, start the gateway:

```bash
make run
```

## 📊 Administration & Observability

- Public API: http://localhost:8080
- Admin API: http://localhost:8081
- Prometheus Metrics: http://localhost:8081/metrics

### Reloading Configuration

To reload the `config.yaml` without stopping the process:

```bash
curl -X POST http://localhost:8081/reload
# OR send a system signal
kill -HUP <PID>
```
