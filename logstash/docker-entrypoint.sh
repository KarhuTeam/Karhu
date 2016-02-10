#!/bin/bash

# Generate Logstash SSL if needed
if [ ! -f "$LOGSTASH_TLS_KEY" ]; then


    mkdir -p `dirname $LOGSTASH_TLS_KEY` || exit 1
    openssl req -subj "/CN=$LOGSTASH_HOST/" -x509 -days 3650 -batch -nodes -newkey rsa:2048 -keyout $LOGSTASH_TLS_KEY -out $LOGSTASH_TLS_CRT || exit 1
else
    echo "Use exising logstash certificats $LOGSTASH_TLS_KEY"
fi

chsum1=""

while true
do
    chsum2=`find /etc/logstash/conf.d -type f -exec md5sum {} \;`
    if [ "$chsum1" != "$chsum2" ]; then

        echo "Config changed, restart logstash"
        /opt/logstash/bin/logstash -f /etc/logstash/conf.d -t && \
            service logstash restart && \
        chsum1=$chsum2
    fi
    sleep 30
done
