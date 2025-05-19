# Start from the official Go image
FROM golang:1.24.2-alpine

RUN apk add --no-cache git ca-certificates

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 
    # GOPROXY=direct

# Create app directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Command to run the executable
CMD ["./main"]
