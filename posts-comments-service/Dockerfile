FROM golang:alpine AS builder

WORKDIR /posts-comments-service
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY posts-comments-service/ ./posts-comments-service/
COPY .env ./
COPY protos/ ./protos/
RUN go build -o service ./posts-comments-service/cmd/main.go

CMD ["./service"]