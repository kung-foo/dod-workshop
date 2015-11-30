FROM smebberson/alpine-consul-redis
MAINTAINER github@jonathan.camp

RUN apk --update add curl

COPY redis-local.conf /etc/
ADD root /
