# Stage 1: Builder
FROM golang:1.23.4-alpine3.21 AS builder

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
RUN go build -o main main.go

# Stage 2: Final Image
FROM alpine:latest

# Install wkhtmltopdf and its dependencies
RUN apk add --no-cache \
    wkhtmltopdf \
    # Required dependencies for wkhtmltopdf
    qt5-qtbase-dev \
    qt5-qtwebkit-dev \
    qt5-qtsvg-dev \
    ttf-dejavu \
    ttf-liberation \
    fontconfig \
    dbus

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/dev.env app.env

# Copy the entrypoint script
COPY entrypoint.sh /app/entrypoint.sh

# Ensure the entrypoint script is executable
RUN chmod +x /app/entrypoint.sh

# Expose the desired port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]