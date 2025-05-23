FROM golang:alpine AS builder

WORKDIR /statistics-service
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY statistics-service/ ./statistics-service/
COPY .env ./
COPY protos/ ./protos/
RUN go build -o service ./statistics-service/cmd/main.go

CMD ["./service"]