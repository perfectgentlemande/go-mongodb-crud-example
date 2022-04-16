FROM golang:1.18.0 AS builder

ADD . /app
WORKDIR /app
# GOOS/GOARCH as you build not from go alpine
RUN GOOS=linux GOARCH=amd64 go build -o go-mongodb-app ./cmd/go-mongodb-crud-api-example

FROM alpine:3.15 AS app
WORKDIR /app
COPY --from=builder /app/go-mongodb-app /app
COPY --from=builder /app/cmd/go-mongodb-crud-api-example/config.yaml /app
CMD ["./go-mongodb-app"]