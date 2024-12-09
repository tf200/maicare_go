FROM golang:1.23.4-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


FROM alpine:latest
# Description

WORKDIR /app
COPY --from=builder /app/main .


EXPOSE 8080

CMD ["/app/main"]

