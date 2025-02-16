FROM golang:1.23.4-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .  

RUN go build -o main ./cmd

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

RUN chmod +x /app/main

ENV DB_SOURCE=postgresql://user:password@postgres:5432/dbname?sslmode=disable

EXPOSE 8081

CMD ["./main"]
