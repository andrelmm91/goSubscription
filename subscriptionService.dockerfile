# # Build stage
# FROM golang:1.17-alpine AS builder

# # Set the current working directory inside the container
# WORKDIR /app

# # Copy the source code into the container
# COPY . .

# # Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/web

# Runtime stage
FROM alpine:latest

RUN mkdir /app

# # Set the current working directory inside the container
# WORKDIR /app

# Copy the binary from the build stage to the runtime image
COPY subscriptionService /app

# Command to run the executable
CMD ["/app/subscriptionService"]
