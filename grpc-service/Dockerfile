FROM golang:1.22-alpine

WORKDIR /app

# Add necessary build tools
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
WORKDIR /app/grpc-service
RUN go build -o main cmd/main.go

# Expose gRPC port
EXPOSE 50051

# Run the application
CMD ["./main"]