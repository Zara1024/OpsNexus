# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /opsnexus-server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /opsnexus-server /usr/local/bin/opsnexus-server
COPY config.yaml /app/config.yaml

USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/opsnexus-server"]
