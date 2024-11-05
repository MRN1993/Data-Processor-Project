# Stage 1: Build
FROM golang:1.22.6 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# If needed, set environment variables
# ENV PORT=8080

# Expose the port if your application uses one
# EXPOSE 8080

# Run the application
CMD ["./main"]
