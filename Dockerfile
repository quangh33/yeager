# --- Stage 1: Builder ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o collector ./cmd/collector/main.go

# --- Stage 2: Runtime ---
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/collector .

EXPOSE 8080

CMD ["./collector"]