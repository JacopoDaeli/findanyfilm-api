FROM golang

RUN go get github.com/tools/godep

EXPORT GIN_MODE=release

ADD . /go/src/github.com/JacopoDaeli/findanyfilm-api
RUN go install findanyfilm-api
ENTRYPOINT /go/bin/findanyfilm-api

EXPOSE 8080
