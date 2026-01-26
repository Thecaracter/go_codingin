# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
# Copy the standalone go.mod
COPY go.mod ./
# Copy source
COPY main.go ./

# Build binary
RUN go build -o webhook-service main.go

# Runtime Stage
FROM alpine:latest

# Deployment script expects this path
WORKDIR /DATA/AppData/backend_codingin

# Install dependencies required for deployment script
RUN apk add --no-cache git docker-cli docker-compose bash curl

# Copy binary from builder
COPY --from=builder /app/webhook-service /usr/local/bin/webhook-service

# Expose port
EXPOSE 9090

# Command
CMD ["webhook-service"]
