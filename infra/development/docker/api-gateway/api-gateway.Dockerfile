FROM alpine
LABEL authors="abbass"

WORKDIR /app

ADD build build

ENTRYPOINT  build/api-gateway