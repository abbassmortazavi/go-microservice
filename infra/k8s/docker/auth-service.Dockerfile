# Build stage - استفاده از alpine base
FROM golang:1.25.5-alpine AS builder

WORKDIR /app


# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./services/auth-service/cmd/

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

EXPOSE 9092

CMD ["/app/auth-service"]