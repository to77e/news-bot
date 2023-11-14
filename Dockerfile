FROM golang:1.21-alpine3.18 AS builder

WORKDIR /app

COPY . .

RUN apk add git

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM scratch as server

WORKDIR /app

COPY --from=builder /app/main ./main

CMD ["./main"]
