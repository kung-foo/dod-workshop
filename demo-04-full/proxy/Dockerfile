FROM smebberson/alpine-consul-base
MAINTAINER github@jonathan.camp

RUN apk --update add curl xz

RUN cd /usr/local/bin && \
    curl -OL https://github.com/kung-foo/waitforit/releases/download/v0.0.1/waitforit-linux-amd64 && \
    chmod +x waitforit-linux-amd64 && \
    curl -L https://github.com/kung-foo/dod-workshop/releases/download/v0.0.1/rrproxy.xz | xz -d > rrproxy && \
    chmod +x rrproxy

RUN cd /tmp && \
    curl -OL https://releases.hashicorp.com/consul-template/0.11.1/consul-template_0.11.1_linux_amd64.zip && \
    unzip consul-template*.zip -d /usr/local/bin && \
    rm -f consul-template*.zip

ADD root /
ADD config.toml.ctmpl /root/
ADD config.toml /root/
