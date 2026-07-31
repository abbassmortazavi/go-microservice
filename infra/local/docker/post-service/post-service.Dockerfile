FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app


COPY build/post-service /app/

EXPOSE 9094

CMD ["./post-service"]