FROM golang:1.21-alpine3.18 AS builder

WORKDIR /news-fetching-bot

COPY . .

RUN apk add git

RUN CGO_ENABLED=0 GOOS=linux go build -o news-fetching-bot ./cmd/news-fetching-bot/main.go

FROM scratch as server

WORKDIR /news-fetching-bot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /news-fetching-bot/news-fetching-bot ./news-fetching-bot

CMD ["./news-fetching-bot"]
