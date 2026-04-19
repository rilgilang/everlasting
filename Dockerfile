FROM golang:1.26.2

ENV GO111MODULE=on

COPY . /app
WORKDIR /app

RUN go build -o ./bin/app