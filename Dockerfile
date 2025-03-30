FROM golang:1.23-alpine AS builder

RUN apk update && apk add --no-cache \
    pkgconfig \
    openssl-dev \
    ca-certificates \
    build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o mattermost-bot ./cmd/

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/mattermost-bot /app/mattermost-bot

EXPOSE 8080
ENTRYPOINT ["/app/mattermost-bot"]