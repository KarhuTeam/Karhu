FROM ubuntu:latest
MAINTAINER max@wayt.me

# Golang path
RUN mkdir -p /go

ENV GO_VERSION 1.5.3
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV WORKDIR ${GOPATH}/src/github.com/karhuteam/karhu
ENV PATH ${PATH}:${GOROOT}/bin
ENV LOGSTASH_TLS_CRT=/etc/logstash/certs/logstash.crt
ENV LOGSTASH_AUTHFILE=/etc/logstash/certs/authfile
ENV LOGSTASH_TAGS_FILTERS=/etc/logstash/conf.d/10-tags-filters.conf
ENV LOGSTASH_APPS_FILTERS=/etc/logstash/conf.d/11-apps-filters.conf
ENV GRAFANA_URL http://localhost:3000

# custom ppa for ansible
RUN apt-get update
RUN apt-get install -y software-properties-common wget apt-transport-https
RUN apt-add-repository ppa:ansible/ansible

# Grafana deps
RUN echo "deb https://packagecloud.io/grafana/stable/debian/ wheezy main" > /etc/apt/sources.list.d/grafana.list
RUN wget -qO - https://packagecloud.io/gpg.key | apt-key add -

# Install ansible && deps
RUN apt-get update && \
    apt-get install -y ansible git grafana && \
    rm -rf /var/lib/apt/lists/*

# Install Golang
RUN wget https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz
RUN tar -C /usr/local -xzf /tmp/go.tar.gz


# Install sources
RUN mkdir -p ${WORKDIR}
ADD . ${WORKDIR}
WORKDIR ${WORKDIR}

RUN go get -v && \
    go build

# Install Grafana
ADD grafana.ini /etc/grafana/grafana.ini

# Setup Logstash
# Default logstash cert path
RUN mkdir -p /logstash
ADD logstash/ /etc/logstash/conf.d/

# Default data path volume
VOLUME /data
VOLUME /logstash

EXPOSE 8080
EXPOSE 3000

ENTRYPOINT ["./docker-entrypoint.sh"]
