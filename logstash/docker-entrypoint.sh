#!/bin/bash

# Generate Logstash SSL if needed
if [ ! -f "$LOGSTASH_TLS_KEY" ]; then


    mkdir -p `dirname $LOGSTASH_TLS_KEY` || exit 1
    openssl req -subj "/CN=karhu/" -x509 -days 3650 -batch -nodes -newkey rsa:2048 -keyout $LOGSTASH_TLS_KEY -out $LOGSTASH_TLS_CRT || exit 1
else
    echo "Use exising logstash certificats $LOGSTASH_TLS_KEY"
fi

code=`curl --write-out "%{http_code}\n" --silent --output /dev/null http://192.168.99.100:9200/_template/filebeat?pretty`

# Generate filebeat index in ES
if [ "$code" != "200" ]; then
    echo "Setup elasticsearch filebeat index template"
    curl --output /tmp/filebeat-index-template.json https://gist.githubusercontent.com/thisismitch/3429023e8438cc25b86c/raw/d8c479e2a1adcea8b1fe86570e42abab0f10f364/filebeat-index-template.json || exit 1
    curl -XPUT 'http://elasticsearch:9200/_template/filebeat?pretty' -d@/tmp/filebeat-index-template.json || exit 1
fi

# Generate collectd authfile
if [ ! -f "$LOGSTASH_COLLECTD_AUTHFILE" ]; then

    echo "Setup collectd authfile"
    echo "collectd: $(openssl rand -base64 32 | md5sum | cut -f1 -d' ')" > $LOGSTASH_COLLECTD_AUTHFILE
fi

chsum1=""

while true
do
    chsum2=`find /etc/logstash/conf.d -type f -exec md5sum {} \;`
    if [ "$chsum1" != "$chsum2" ]; then

        date=$(date)
        echo "$date Config changed, restart logstash"
        /opt/logstash/bin/logstash -f /etc/logstash/conf.d -t && \
            service logstash restart && \
        chsum1=$chsum2
    fi
    sleep 30
done
