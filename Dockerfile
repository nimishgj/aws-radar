# Build stage
FROM golang:alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/

# Build with optimizations for smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o aws-radar ./cmd/aws-radar

# Runtime stage - minimal alpine for utilities
FROM alpine:3.21

# Install CA certificates, timezone data, and wget for healthcheck
RUN apk add --no-cache ca-certificates tzdata wget

# Copy binary
COPY --from=builder /build/aws-radar /app/aws-radar

# Copy default config
COPY config.yaml /app/config.yaml

WORKDIR /app

# Expose metrics port
EXPOSE 9090

# Run the binary
ENTRYPOINT ["./aws-radar"]
