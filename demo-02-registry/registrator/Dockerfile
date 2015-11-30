FROM smebberson/alpine-base:1.2.0
MAINTAINER github@jonathan.camp

RUN apk --update add curl go git mercurial

RUN cd /usr/local/bin && \
    curl -OL https://github.com/kung-foo/waitforit/releases/download/v0.0.1/waitforit-linux-amd64 && \
    chmod +x waitforit-linux-amd64

RUN mkdir -p /go/src/github.com/gliderlabs/ \
    && cd /go/src/github.com/gliderlabs/ \
    && git clone -b v6 --single-branch https://github.com/gliderlabs/registrator \
    && cd registrator \
    && export GOPATH=/go \
    && go get \
    && go build -ldflags "-X main.Version $(cat VERSION)" -o /bin/registrator

ADD root /
