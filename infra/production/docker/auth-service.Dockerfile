# Build stage
FROM golang:1.25.5 AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files (go.sum may not exist for new projects)
COPY go.mod ./

# Check if go.sum exists, if not run go mod tidy
RUN if [ ! -f go.sum ]; then go mod tidy; fi

# Copy source code
COPY . .

# Build the application
# Adjust the path based on where your main.go is located
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./services/auth-service/cmd/

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder --chown=app:app /app/auth-service /app/auth-service


# Switch to non-root user
USER app

EXPOSE 9091

CMD ["/app/auth-service"]