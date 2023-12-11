FROM golang:1.20-alpine
LABEL authors="andreisarnouski"

WORKDIR /app

RUN apk add --no-cache make
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh

RUN rm -f $GOPATH/go.mod $GOPATH/go.sum

COPY . .

RUN  go mod download

RUN make check
#RUN make gooseUp

COPY .env cmd/consumer
COPY .env internal/kafka/

WORKDIR /app/cmd/consumer

RUN go build -o app .

CMD ["./app"]