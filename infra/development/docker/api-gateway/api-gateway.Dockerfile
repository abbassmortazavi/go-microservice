# Build stage
FROM golang:1.25.5 AS builder

WORKDIR /app

# Copy go mod files (go.sum may not exist for new projects)
COPY go.mod ./

# Check if go.sum exists, if not run go mod tidy
RUN if [ ! -f go.sum ]; then go mod tidy; fi

# Copy source code
COPY . .

# Build the application
# Adjust the path based on where your main.go is located
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-gateway ./services/api-gateway

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway /app/

# Run the binary
ENTRYPOINT ["/app/api-gateway"]