FROM golang:alpine AS builder

WORKDIR /api-gateway
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY api-gateway/ ./api-gateway/
COPY .env ./
RUN go build -o gateway ./api-gateway/cmd/main.go

CMD ["./gateway"]