FROM golang:1.21-alpine
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

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache \
    go mod download

COPY .env cmd/consumer
#COPY .env internal/kafka/

WORKDIR /app/cmd/consumer

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go build -o app .

CMD ["./app"]