# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set working directory in container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go app (output binary as 'main')
RUN go build -o main .

# Stage 2: Create a smaller, production-ready image
FROM alpine:latest

# Install necessary dependencies to run the Go app
RUN apk --no-cache add ca-certificates

# Set working directory in container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 to access the app
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
