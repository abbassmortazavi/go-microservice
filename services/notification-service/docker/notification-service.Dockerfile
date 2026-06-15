FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY build/notification-service /app/
EXPOSE 9093
CMD ["./notification-service"]