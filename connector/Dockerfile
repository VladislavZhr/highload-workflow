# syntax=docker/dockerfile:1

FROM golang:1.26.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin/connector \
    ./cmd/app

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/connector /app/connector

EXPOSE 8080

ENTRYPOINT ["/app/connector"]