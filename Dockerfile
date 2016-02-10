FROM ubuntu:latest
MAINTAINER max@wayt.me

# Golang path
RUN mkdir -p /go

ENV GO_VERSION 1.5.3
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV WORKDIR ${GOPATH}/src/github.com/karhuteam/karhu
ENV PATH ${PATH}:${GOROOT}/bin
ENV LOGSTASH_TLS_KEY /logstash/logstash.key
ENV LOGSTASH_TLS_CRT /logstash/logstash.crt

# custom ppa for ansible
RUN apt-get update
RUN apt-get install -y software-properties-common wget
RUN apt-add-repository ppa:ansible/ansible

# Logstash deps
RUN wget -qO - https://packages.elastic.co/GPG-KEY-elasticsearch | apt-key add -
RUN echo 'deb http://packages.elastic.co/logstash/2.2/debian stable main' | tee /etc/apt/sources.list.d/logstash-2.2.x.list

# Install ansible && deps
RUN apt-get update && \
    apt-get install -y ansible git logstash

# Install Golang
RUN wget https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz
RUN tar -C /usr/local -xzf /tmp/go.tar.gz


# Install sources
RUN mkdir -p ${WORKDIR}
ADD . ${WORKDIR}
WORKDIR ${WORKDIR}

RUN go get -v && \
    go build

# Setup Logstash
# Default logstash cert path
RUN mkdir -p /logstash
ADD logstash/ /etc/logstash/conf.d/

# Default data path volume
VOLUME /data
VOLUME /logstash

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
