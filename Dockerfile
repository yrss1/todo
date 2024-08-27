FROM golang:1.22.5-alpine as builder

WORKDIR /build

COPY . /build

RUN go mod download

COPY . .

RUN go build -o todo-service .

FROM alpine:3.18 as hoster
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/todo-service ./todo-service
COPY --from=builder /build/db/migrations ./db/migrations

ENTRYPOINT ["./todo-service"]