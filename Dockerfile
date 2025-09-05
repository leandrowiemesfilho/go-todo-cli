# Build stage
FROM golang:1.24.6-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o go-todo-cli ./cmd/main.go

# Final stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/go-todo-cli .
COPY --from=builder /app/migrations ./migrations

RUN chmod +x go-todo-cli

# Create non-root user
RUN adduser -D -u 1000 todouser && \
    chown -R todouser:todouser /app

# Switch to non-root user
USER todouser

# Command to run the application
ENTRYPOINT ["./go-todo-cli"]