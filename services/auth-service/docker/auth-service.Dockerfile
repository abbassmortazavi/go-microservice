# Build stage
FROM golang:1.25.5 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application - مسیر صحیح
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./services/auth-service

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/auth-service .

# Run the binary
ENTRYPOINT ["/app/auth-service"]