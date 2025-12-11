# Build stage
FROM golang:1.25.5-alpine AS builder
WORKDIR /app

# Install air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Install swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Development stage (with air for hot-reload)
FROM builder AS dev
CMD ["air"]

