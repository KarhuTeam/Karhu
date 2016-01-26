FROM ubuntu:latest
MAINTAINER max@wayt.me

# Golang path
RUN mkdir -p /go

ENV GO_VERSION 1.5.3
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV WORKDIR ${GOPATH}/src/github.com/karhuteam/karhu
ENV PATH ${PATH}:${GOROOT}/bin

# custom ppa for ansible
RUN apt-get update
RUN apt-get install -y software-properties-common
RUN apt-add-repository ppa:ansible/ansible

# Install ansible && deps
RUN apt-get update && \
    apt-get install -y ansible wget git

# Install Golang
RUN wget https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz
RUN tar -C /usr/local -xzf /tmp/go.tar.gz


# Install sources
RUN mkdir -p ${WORKDIR}
ADD . ${WORKDIR}
WORKDIR ${WORKDIR}

RUN go get -v && \
    go build

# Default data path volume
VOLUME /data

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
