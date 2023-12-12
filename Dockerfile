FROM golang:1.20-alpine
LABEL authors="andreisarnouski"

WORKDIR /app

RUN apk add --no-cache make
RUN apk add --no-cache bash
RUN apk add nano
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh

COPY . .

RUN  go mod download

COPY .env cmd/consumer
COPY .env internal/kafka/

WORKDIR /app/cmd/consumer

RUN go build -o app .

CMD ["./app"]