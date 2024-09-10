# Stage 1: Build the Go binary
FROM golang:1.18-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Stage 2: Run the Go binary
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder     /app/main .
COPY static/config.json /app/
COPY ./privatekey.key   /root/
COPY ./publiccert.crt   /root/
# shared volume to exchange nginx configs
RUN  mkdir -p /root/nginx
# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main", "/app/config.json"]