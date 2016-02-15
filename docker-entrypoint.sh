#!/bin/bash

# Generate Logstash SSL if needed
# if [ ! -f "$LOGSTASH_TLS_KEY" ]; then
#
#
#     mkdir -p `dirname $LOGSTASH_TLS_KEY` || exit 1
#     openssl req -subj "/CN=$LOGSTASH_HOST/" -x509 -days 3650 -batch -nodes -newkey rsa:2048 -keyout $LOGSTASH_TLS_KEY -out $LOGSTASH_TLS_CRT || exit 1
# else
#     echo "Use exising logstash certificats $LOGSTASH_TLS_KEY"
# fi

echo "Ensure data dir exist"
mkdir -p $STORAGE_PATH

# echo "DEBUG=${DEBUG:-1}
# MGO_HOSTS=${MGO_HOSTS:-mongo}
# MGO_DB=${MGO_DB:-karhu}
#
# PUBLIC_HOST=${PUBLIC_HOST:-http://127.0.0.1:8080}
#
# STORAGE_DRIVER=${STORAGE_DRIVER:-fs}
# STORAGE_PATH=$KARHU_STORAGE_PATH" > .env
#
# cat .env; echo

service grafana start

./karhu
