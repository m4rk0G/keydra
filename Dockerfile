# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o bin/keydra \
    ./cmd/keydra

# Final stage
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 keydra && \
    adduser -D -s /bin/sh -u 1000 -G keydra keydra

WORKDIR /

# Copy binary from builder stage
COPY --from=builder /app/bin/keydra /usr/local/bin/keydra

# Change ownership and permissions
RUN chown keydra:keydra /usr/local/bin/keydra && \
    chmod +x /usr/local/bin/keydra

# Switch to non-root user
USER keydra

# Expose metrics port (optional)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD pgrep keydra || exit 1

# Default command
ENTRYPOINT ["/usr/local/bin/keydra"]
CMD ["-config=/etc/keydra/config.yaml"]
