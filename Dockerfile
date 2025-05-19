# Start from the official Go image
FROM golang:1.24.2

RUN apt-get update && apt-get install -y \
    git \
    ca-certificates \
    curl \
    build-essential \
    librdkafka-dev


# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 

# Create app directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN go build -o main ./cmd

# Command to run the executable
CMD ["./main"]
