FROM golang

RUN go get github.com/tools/godep

ENV GIN_MODE=release

ADD . /go/src/findanyfilm-api
RUN go install findanyfilm-api
ENTRYPOINT /go/bin/findanyfilm-api

EXPOSE 8080
