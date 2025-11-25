# Stage 1: Build air and goose in a builder stage
FROM golang:1.25-alpine AS builder

# Install the build tools
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Stage 2: Final minimal runtime stage
FROM golang:1.25-alpine AS runtime

# Install only runtime dependencies (no build tools needed)
RUN apk add --no-cache \
    git \
    make \
    curl \
    ffmpeg \
    ca-certificates

# Copy pre-built binaries from builder stage
COPY --from=builder /go/bin/air /go/bin/air
COPY --from=builder /go/bin/goose /go/bin/goose

# Set working directory
WORKDIR /app

# Fix git dubious ownership error for mounted volumes
RUN git config --global --add safe.directory /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Expose application port
EXPOSE 60000

# Start air (uses default .air.toml config)
CMD ["air"]
