FROM ubuntu:latest
MAINTAINER max@wayt.me

ENV LOGSTASH_TLS_KEY /etc/logstash/certs/logstash.key
ENV LOGSTASH_TLS_CRT /etc/logstash/certs/logstash.crt
ENV LOGSTASH_COLLECTD_AUTHFILE /etc/logstash/certs/authfile

# custom ppa
RUN apt-get update
RUN apt-get install -y software-properties-common wget curl
RUN add-apt-repository -y ppa:webupd8team/java

# Logstash deps
RUN wget -qO - https://packages.elastic.co/GPG-KEY-elasticsearch | apt-key add -
RUN echo 'deb http://packages.elastic.co/logstash/2.2/debian stable main' | tee /etc/apt/sources.list.d/logstash-2.2.x.list

# Accept java licence
RUN echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | debconf-set-selections

# Install ansible && deps
RUN apt-get update && \
    apt-get install -y logstash oracle-java8-installer && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /var/cache/oracle-jdk8-installer

# Setup Logstash
# Default logstash cert path
RUN mkdir -p /etc/logstash/certs && \
    mkdir -p /etc/logstash/conf.d
ADD docker-entrypoint.sh /
ADD conf.d/02-beats-input.conf /etc/logstash/conf.d
ADD conf.d/03-collectd-input.conf /etc/logstash/conf.d
ADD conf.d/10-syslog-filter.conf /etc/logstash/conf.d
ADD conf.d/30-elasticsearch-output.conf /etc/logstash/conf.d

# Default data path volume
VOLUME /etc/logstash

WORKDIR /

EXPOSE 5044
EXPOSE 25826

ENTRYPOINT ["./docker-entrypoint.sh"]
