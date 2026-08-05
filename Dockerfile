# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

# Install git and other necessary tools
# (Required if fetching dependencies from private repos or specific CGO requirements)
RUN apk update && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached unless go.mod changes)
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# CGO_ENABLED=0 ensures a statically linked binary which can run in the scratch/alpine image
# GOOS=linux ensures we are building for linux, regardless of the host OS
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o my_app .

# Stage 2: Create a minimal production image
FROM alpine:latest

# Install CA certificates for making secure HTTPS requests (often needed in backend services)
# Also add tzdata if timezones are important
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/my_app .

# Expose the port your application listens on
EXPOSE 8081

# Command to run the executable
CMD ["./my_app"]
