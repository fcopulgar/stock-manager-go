# Stage 1: Build the Go binary
FROM golang:1.23.3-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main cmd/main.go

# Stage 2: Run the binary
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

# Expose port (if your app uses any network port)
# EXPOSE 8080

# Command to run the executable
CMD ["./main"]
