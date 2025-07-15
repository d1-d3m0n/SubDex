# Use official Go image as base
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN go build -o subdex .

# Run the tool by default
ENTRYPOINT ["./subenum"]
