FROM library/golang:1.5-alpine
MAINTAINER github@jonathan.camp

RUN apk --update add git make
RUN go get github.com/tools/godep

COPY . /usr/src/app

RUN cd /usr/src/app && make
