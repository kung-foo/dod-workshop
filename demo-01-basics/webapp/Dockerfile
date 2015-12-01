FROM smebberson/alpine-base:1.2.0
MAINTAINER github@jonathan.camp

RUN     echo "@testing http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN     apk add --update curl python py-flask py-netifaces@testing && rm -rf /var/cache/apk/*
WORKDIR /root
ADD     app.py /root/
CMD     ["python", "app.py"]
