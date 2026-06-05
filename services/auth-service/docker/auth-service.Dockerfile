FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app


COPY build/auth-service /app/

EXPOSE 9092

CMD ["./auth-service"]