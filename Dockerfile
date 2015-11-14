FROM golang:1.5

ADD . /go/src/github.com/jmank88/todo
WORKDIR /go/src/github.com/jmank88/todo

ENV GO15VENDOREXPERIMENT 1
RUN go build ./...

EXPOSE 8080