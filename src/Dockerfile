# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# add gcc
RUN apk add --no-cache gcc musl-dev
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest


# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o main .

# Stage 2: Run the Go binary
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder     /go/bin/dlv /
COPY --from=builder     /app/main .
COPY static/config.json /app/
COPY ./privatekey.key   /root/
COPY ./publiccert.crt   /root/
# shared volume to exchange nginx configs
RUN  mkdir -p /root/nginx
# Expose port 8080 to the outside world, 40000 is debugger
EXPOSE 8080 40000

# Command to run the executable
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main", "/app/config.json" ]