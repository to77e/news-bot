FROM golang:1.21-alpine3.18 as builder

RUN apk add git

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /migrations

COPY migrations/ /migrations/

FROM alpine:3.18

COPY --from=builder /go/bin/goose /go/bin/goose
COPY --from=builder /migrations/ /migrations/

ENV DATABASE_HOST=""
ENV DATABASE_PORT=""
ENV DATABASE_USER=""
ENV DATABASE_PASSWORD=""
ENV DATABASE_NAME=""

ENTRYPOINT ["/bin/sh", "-c", "/go/bin/goose -dir /migrations postgres \"host=${DATABASE_HOST} port=${DATABASE_PORT} user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} sslmode=disable\" up"]