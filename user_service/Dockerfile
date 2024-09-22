# Stage 1: Build the Go binary using Go 1.23
FROM golang:1.23-alpine AS builder

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application binary
RUN go build -o main .

# Stage 2: Create a lightweight image with only the binary
FROM alpine:latest

# Set the working directory for the final image
WORKDIR /root/

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port on which the application will run (if needed)
EXPOSE 8080

# Run the Go binary
CMD ["./main", "rest"]
