# FROM golang:1.18.0-alpine3.15 AS builder

# ADD . /app
# WORKDIR /app
# RUN go mod tidy
# RUN go build -o go-mongodb-app ./cmd/go-mongodb-crud-api-example

FROM golang:1.18.0 AS builder

ADD . /app
WORKDIR /app
RUN go build -o go-mongodb-app ./cmd/go-mongodb-crud-api-example

FROM alpine:3.15 AS app
WORKDIR /app
COPY --from=builder /app/go-mongodb-app /app
CMD ["./go-mongodb-app"]