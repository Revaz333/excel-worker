FROM golang:latest

ARG APP_DIR=app

COPY . /go/tmp/src/${APP_NAME}

WORKDIR /go/tmp/src/${APP_NAME}

RUN go mod tidy
#RUN go build ./cmd/main.go
