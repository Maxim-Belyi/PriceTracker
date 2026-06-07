FROM golang:1.26-alpine AS builder
 
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o worker_app cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/worker_app .

CMD ["./worker_app"]