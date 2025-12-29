# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]

