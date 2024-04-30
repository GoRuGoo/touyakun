FROM golang:1.22.2-alpine

RUN apk update && apk add git

RUN mkdir /go/src/api

WORKDIR /go/src/api

ADD . /go/src/api

RUN go install github.com/cosmtrek/air@latest

RUN touch .air.toml

EXPOSE 8080

CMD [ "air","-c",".air.toml" ]