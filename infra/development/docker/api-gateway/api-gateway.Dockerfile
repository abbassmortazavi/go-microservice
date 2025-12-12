FROM alpine
LABEL authors="abbass"

WORKDIR /app

ADD build build

ENTRYPOINT exec build/api-gateway