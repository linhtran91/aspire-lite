# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yml config/config.yml

EXPOSE 8080
CMD [ "/cmd/main" ]
