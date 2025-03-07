FROM golang:alpine AS builder

WORKDIR /user-service
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY user-service/ ./user-service/
COPY .env ./
RUN go build -o service ./user-service/cmd/main.go

CMD ["./service"]