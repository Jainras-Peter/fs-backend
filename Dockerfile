# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS build
WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
COPY vendor/ vendor/

# Copy the rest of the source code
COPY . .

# Build a statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o fs-backend .

# Stage 2: Minimal runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the compiled binary and config from the build stage
COPY --from=build /app/fs-backend .

# Expose the server port
EXPOSE 5000

# Start the server
CMD ["./fs-backend"]
