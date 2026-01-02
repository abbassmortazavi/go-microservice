# Build stage - استفاده از alpine base
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Install dependencies (apk برای Alpine کار می‌کند)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY services/auth-service ./services/auth-service
COPY pkg ./pkg
COPY proto ./proto

# Build the application
WORKDIR /app/services/auth-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/auth-service .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -S app -G app

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/auth-service .

# Switch to non-root user
USER app

EXPOSE 9092  # توجه: پورت باید 9092 باشد نه 9091

CMD ["/app/auth-service"]