# Stage 1: Build
FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ngatex ./cmd/gateway

# Stage 2: Final Image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ngatex .
COPY --from=builder /app/config.yaml .
# Run the gateway
ENTRYPOINT ["./ngatex", "--config", "config.yaml"]
EXPOSE 8080 8081