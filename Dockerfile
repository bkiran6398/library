# syntax=docker/dockerfile:1.7

FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/library-http ./cmd/library-http

FROM alpine:latest
WORKDIR /app
COPY --from=builder /out/library-http /app/library-http
COPY config/config.yaml /app/config/config.yaml
COPY migrations /app/migrations
ENV LIB_CONFIG=/etc/library/config.yaml
EXPOSE 8080
ENTRYPOINT ["/app/library-http"]
