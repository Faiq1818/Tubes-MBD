FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
# COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

COPY .env .env

EXPOSE 8000

CMD ["./main"]
