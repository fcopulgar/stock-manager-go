# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install git and other dependencies if necessary
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Set environment variables for a statically linked binary
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go app
RUN go build -o main cmd/main.go
