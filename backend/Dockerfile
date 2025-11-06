# -------- Stage 1: Builder --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for go mod) — no Delve this time
RUN apk add --no-cache git

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build optimized binary
RUN go build -o url-shortener ./cmd/server

# -------- Stage 2: Runtime --------
FROM alpine:latest

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/url-shortener .
COPY config.yaml .
COPY migrations ./migrations

# Expose server port
EXPOSE 8080

# Run the app normally (no debugger)
CMD ["./url-shortener"]
