# Stage 1: Builder
FROM golang:1.24.9-alpine3.22 AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files first to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

ENV GOPROXY=https://proxy.golang.org,direct

# Build the Go application
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main main.go

# Stage 2: Final Image
FROM alpine:3.22

# Install runtime certificates for HTTPS integrations
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

COPY db/migrations /app/db/migrations

# Expose the desired port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/app/main"]
