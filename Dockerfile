# Stage 1: Build Module
FROM golang:1.25-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Update and install git (in case dependencies need it)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build
# -o main: name the binary "main"
# ./cmd/server: path to main package
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

# Security: Add a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Install certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port (Documentation purpose, actual mapping in docker-compose)
EXPOSE 8080

# Run the binary
CMD ["./main"]
