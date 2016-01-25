#!/bin/bash

KARHU_STORAGE_PATH=${STORAGE_PATH:-/data}

echo "Ensure data dir exist"
mkdir -p $KARHU_STORAGE_PATH


echo "DEBUG=${DEBUG:-1}
MGO_HOSTS=${MGO_HOSTS:-mongo}
MGO_DB=${MGO_DB:-karhu}

PUBLIC_HOST=${PUBLIC_HOST:-http://127.0.0.1:8080}

STORAGE_DRIVER=${STORAGE_DRIVER:-fs}
STORAGE_PATH=$KARHU_STORAGE_PATH" > .env

cat .env; echo

./karhu
