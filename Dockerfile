# base image
FROM golang:1.22.0 AS builder

# Declaring env variables
ENV GOOS=linux \
    CGO_ENABLED=0

WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build
COPY . .
RUN go build -o profile-builder ./cmd/main.go

# Final stage
FROM alpine:latest AS production

# Add ca-certificates for HTTPS calls
RUN apk add --no-cache ca-certificates

# Set work directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/profile-builder /app/profile-builder

# Command to run
CMD ["/app/profile-builder"]

# Expose port
EXPOSE 8080