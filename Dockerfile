# -------- Stage 1: Builder --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o url-shortener ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener .
COPY config.yaml .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./url-shortener"]
