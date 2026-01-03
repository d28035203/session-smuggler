# Start from Golang base image
FROM golang:1.20-alpine AS builder

# Install necessary dependencies
RUN apk update && apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o app .

# Build a small image from scratch
FROM alpine:3.16

# Set the working directory in the final image
WORKDIR /app

# Copy the compiled binary from the builder image
COPY --from=builder /app/app .

# Expose the port the app runs on
EXPOSE 8085

# Command to run the executable
CMD ["./app"]