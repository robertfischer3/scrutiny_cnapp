FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o build/myapp ./cmd/myapp

# Start a new stage from scratch
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/build/scrutiny_cnapp .
COPY --from=builder /app/configs ./configs

# Expose application port
EXPOSE 8080

# Command to run the executable
CMD ["./scrutiny_cnapp"]
