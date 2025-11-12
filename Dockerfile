# syntax=docker/dockerfile:1.6
# Multi-stage build for smaller final image
FROM golang:1.24.3-alpine AS builder
 
# Install git for go mod dependencies
#apk is the package manager for alpine linux like npm for node
# we want no cache because we wont need to install anything else after the initial build
RUN apk add --no-cache git gcc musl-dev

# Set working directory
#we are creating a directory called app in the root of the project
WORKDIR /app

# Copy go mod files first for better caching
# were copying go.mod and go.sum to the base directory we created above
COPY go.mod go.sum ./

# Download dependencies
# this will download all the dependencies for the project
RUN --mount=type=cache,target=/go/pkg/mod \
    GOPROXY=direct \
    GOSUMDB=off \
    go mod download -x

# Copy source code
# this will copy all the files in the root of the project to the working directory we created above
COPY . .

# Build the application from the correct path
# CGO_ENABLED=0 means we are disabling the ability to reference any C libraries
# GOOS=linux means we are building the binary for linux
# -a means we are building the binary for the current platform
# -installsuffix cgo allows us to separate the cache for the binary so we dont use a different cache or any preexisting cachedfiles
# -o is output the binary to the main 
# last one points to the main.go file in the src/api/v1 directory
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    go build -v -p 2 \
      -buildvcs=false -trimpath -ldflags "-s -w" \
      -o /app/main ./src/api

# Final stage
# this is the final stage of the build
# we are using alpine linux because it is a small linux distribution that is easy to use and has a small footprint
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
# if someone hacks the app they wont have root access to the server
RUN adduser -D -s /bin/sh appuser

# set the working directory to the root of the project
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appuser ./

# Switch to non-root user
USER appuser

# Expose port (default 8080, can be overridden)
EXPOSE 3000

# Command to run
CMD ["./main"]