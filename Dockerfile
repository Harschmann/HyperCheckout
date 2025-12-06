# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# -o main : name the output binary "main"
RUN go build -o main cmd/api/main.go

# Stage 2: Create a tiny runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Copy the .env file (optional, but good for defaults)
COPY --from=builder /app/.env .

# Expose the port
EXPOSE 8080

# Command to run the app
CMD ["./main"]