# Build stage
FROM golang:1.24-bullseye AS builder
RUN apt-get update && apt-get install -y gcc libc6-dev libolm-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# RUN go generate ./...
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Serve stage
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y ca-certificates libolm3 && rm -rf /var/lib/apt/lists/*
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
CMD ["./main", "-debug"]