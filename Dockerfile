# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22.0 AS build-stage
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for version information
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the application with version information
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /api ./cmd/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

# Install ca-certificates for HTTPS and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user to run the app
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

WORKDIR /app

# Copy the binary from build stage
COPY --from=build-stage /api /app/api

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the port (should match PORT env var)
EXPOSE 8080

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Use exec form for proper signal handling
ENTRYPOINT ["/app/api"]