# Build stage
FROM golang:1.24.1-alpine AS builder

# Install git for module downloads
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o processor ./cmd/processor

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/processor .

# Copy sample data
COPY --from=builder /app/data ./data

# Create data directories
RUN mkdir -p /data /data/processed && \
    chown -R appuser:appgroup /app /data

# Switch to non-root user
USER appuser

# Expose port (if needed for health checks)
EXPOSE 8080

# Set default environment variables
ENV ENVIRONMENT=production \
    LOG_LEVEL=info \
    WATCH_DIRECTORY=/data \
    PROCESSED_DIRECTORY=/data/processed

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep processor || exit 1

# Run the application
CMD ["./processor"]