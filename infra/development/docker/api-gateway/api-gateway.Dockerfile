FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# کپی فایل اجرایی
COPY build/api-gateway .


EXPOSE 8081

CMD ["./api-gateway"]