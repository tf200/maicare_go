# Stage 1: Builder
FROM golang:1.23.4-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Stage 2: Final Image
FROM alpine:latest

# Install any necessary packages (optional)
# RUN apk add --no-cache bash

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy the entrypoint script
COPY entrypoint.sh /app/entrypoint.sh

# Ensure the entrypoint script is executable
RUN chmod +x /app/entrypoint.sh

# Expose the desired port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]

# (Optional) If you prefer CMD, you can use:
# CMD ["/app/main"]