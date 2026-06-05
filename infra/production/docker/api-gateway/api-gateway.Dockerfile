FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY build/api-gateway .


EXPOSE 8085

CMD ["./api-gateway"]