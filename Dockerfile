FROM golang:1.21.4

ENV GO111MODULE=on

COPY . /app
WORKDIR /app

RUN go build -o ./bin/app