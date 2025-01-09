FROM golang:1.23-alpine AS builder

# Create a working directory
WORKDIR /app

# Copy Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN go build -o app-binary cmd/main/main.go

# Start a new minimal image
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copy the binary and configs
COPY --from=builder /app/app-binary .
COPY .env /app/.env

# Expose ports
EXPOSE 8080

# Run the binary
CMD ["/app/app-binary"]