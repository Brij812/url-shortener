FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o url-shortener ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/url-shortener .
EXPOSE 8080
CMD ["./url-shortener"]
